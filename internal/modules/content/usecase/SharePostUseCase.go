package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SharePostUseCase struct {
	eventbus          events.EventBus
	postRepo          IRepositoryMongodb.IPostRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postinsightRepo   IRepositoryCassandra.IPostInsights
	postsetting       IRepositoryMongodb.IPostSettingRepository
	pool              IRepositoryShare.IWorkerPool
}

func NewSharePostUseCase(eventbus events.EventBus, postRepo IRepositoryMongodb.IPostRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, pool IRepositoryShare.IWorkerPool) *SharePostUseCase {
	return &SharePostUseCase{
		eventbus:          eventbus,
		postRepo:          postRepo,
		postExtensionRepo: postExtensionRepo,
		pool:              pool,
	}
}
func (s *SharePostUseCase) Execute(ctx context.Context, req *req.SharePostRequest) (*response.Response, error) {
	type taskResult struct {
		datapost *entity.Post
		data     *entity.PostExtension
		err      error
	}
	worker := 2
	resultChan := make(chan taskResult, worker)
	for i := 0; i < worker; i++ {
		err := s.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			switch i {
			case 0:
				datapost, err := s.postRepo.GetPostByID(safectx, req.PostID)
				resultChan <- taskResult{datapost: datapost, err: err}
			case 1:
				postExtension, err := s.postExtensionRepo.GetByPostID(safectx, req.PostID)
				resultChan <- taskResult{data: postExtension, err: err}
			}
		})
		if err != nil {
			return nil, err
		}
	}
	var datapost *entity.Post
	var postExtension *entity.PostExtension
	for i := 0; i < worker; i++ {
		result := <-resultChan
		if result.err != nil {
			return nil, result.err
		}
		if result.datapost != nil {
			datapost = result.datapost
		}
		if result.data != nil {
			postExtension = result.data
		}
	}
	if datapost == nil || postExtension == nil {
		return nil, errors.New("failed to retrieve post data")
	}
	newid := primitive.NewObjectID()
	worker2 := 5
	for i := 0; i < worker2; i++ {
		err := s.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			switch i {
			case 0:
				datapost.Stats.Shares += 1
				_, err := s.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					return
				}
			case 1:
				datapost.UserID = req.UserID

				datapost.ID = newid
				_, err := s.postRepo.CreatePost(safectx, datapost)
				if err != nil {
					return
				}
			case 2:
				postExtension.PostID = newid
				newparent, err := primitive.ObjectIDFromHex(req.PostID)
				if err != nil {
					return
				}
				newidExtension := primitive.NewObjectID()
				postExtension.ID = newidExtension
				postExtension.ShareData.ParentPostID = newparent
				err = s.postExtensionRepo.CreatePostExtension(safectx, postExtension)
				if err != nil {
					return
				}
			case 3:
				idinsight := newid.String()
				postIDUUID, err := gocql.ParseUUID(idinsight)
				if err != nil {
					return
				}
				createPostInsight := &entity.PostInsight{
					PostID:         postIDUUID,
					Reach:          0,
					Impressions:    0,
					EngagementRate: 0,
					ReactionsTotal: 0,
					CommentsTotal:  0,
					SharesTotal:    0,
					ClicksTotal:    0,
					VideoViews3s:   0,
				}
				err = s.postinsightRepo.CreatePostInsightInitPost(safectx, createPostInsight)
				if err != nil {
					return
				}
			case 4:
				createSetting := &entity.PostSetting{
					ID:     primitive.NewObjectID(),
					PostID: newid,
				}
				_, err := s.postsetting.CreatePostSetting(safectx, createSetting)
				if err != nil {
					return
				}
			}
		})
		if err != nil {
			return nil, err
		}
	}
	return response.NewResponse(response.WithData(postExtension),
		response.WithMessage(" "), response.WithStatus(http.StatusOK)), nil
}
