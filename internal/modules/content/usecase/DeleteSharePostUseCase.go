package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type DeleteSharePostUseCase struct {
	eventbus          events.EventBus
	postRepo          IRepositoryMongodb.IPostRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postinsightRepo   IRepositoryCassandra.IPostInsights
	postsetting       IRepositoryMongodb.IPostSettingRepository
	pool              IRepositoryShare.IWorkerPool
}

func NewDeleteSharePostUseCase(eventbus events.EventBus, postRepo IRepositoryMongodb.IPostRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postinsightRepo IRepositoryCassandra.IPostInsights, postsetting IRepositoryMongodb.IPostSettingRepository, pool IRepositoryShare.IWorkerPool) *DeleteSharePostUseCase {
	return &DeleteSharePostUseCase{
		eventbus:          eventbus,
		postRepo:          postRepo,
		postExtensionRepo: postExtensionRepo,
		postinsightRepo:   postinsightRepo,
		postsetting:       postsetting,
		pool:              pool,
	}
}
func (s *DeleteSharePostUseCase) Execute(ctx context.Context, req *contentEvent.SharePostPayload) (*response.Response, error) {
	datapost, err := s.postRepo.GetPostByID(ctx, req.PostID)
	if err != nil {
		return nil, err
	}
	worker := 4
	for i := 0; i < worker; i++ {
		err := s.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			switch i {
			case 0:
				datapost.Stats.Shares -= 1
				_, err = s.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					return
				}
			case 1:
				err = s.postRepo.DeletePost(safectx, req.PostID)
				if err != nil {
					return
				}
			case 2:
				err = s.postExtensionRepo.DeleteByPostID(safectx, req.PostID)
				if err != nil {
					return
				}
			case 3:
				err = s.postsetting.DeletePostSetting(safectx, req.PostID)
				if err != nil {
					return
				}
			case 4:
				err = s.postinsightRepo.DeletePostInsightByPostID(safectx, req.PostID)
				if err != nil {
					return
				}
			}
		})
		log.Printf("Worker %d started\n", i)
		if err != nil {
			log.Printf("Error starting worker %d: %v\n", i, err)
		}
		log.Printf("Worker %d finished\n", i)
	}
	s.pool.Wait()
	return response.NewResponse(response.WithData(""),
		response.WithMessage(" "), response.WithStatus(http.StatusOK)), nil
}
