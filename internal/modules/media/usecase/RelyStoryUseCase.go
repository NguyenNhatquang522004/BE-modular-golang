package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/gocql/gocql"
)

type RelyStoryUseCase struct {
	// Define dependencies for RelyStoryUseCase here (e.g., repositories, services)
	storyRepo IRepositoryMongodb.IStoryRepository
	storyview IRepositoryCassandra.IStoryViewRepository
	pool      IRepositoryShare.IWorkerPool
}

func NewRelyStoryUseCase(
	storyRepo IRepositoryMongodb.IStoryRepository,
	storyview IRepositoryCassandra.IStoryViewRepository,
	pool IRepositoryShare.IWorkerPool,
) *RelyStoryUseCase {
	return &RelyStoryUseCase{
		storyRepo: storyRepo,
		storyview: storyview,
		pool:      pool,
	}
}

func (uc *RelyStoryUseCase) Execute(ctx context.Context, req *req.RelyStoryRequest) ([]*res.FailedStoryResponse, error) {
	// Implement the logic for relying to a story here
	workercount := 2
	resultChan := make(chan *res.FailedStoryResponse, workercount)
	for i := 0; i < workercount; i++ {
		err := uc.pool.Run(ctx, func() {
			var err error
			switch i {
			case 0:
				datastory, err := uc.storyRepo.GetStoryByID(ctx, req.StoryId)
				if err != nil {
					resultChan <- &res.FailedStoryResponse{
						StoryID:      req.StoryId,
						UserID:       req.UserId,
						ErrorMessage: err.Error(),
					}
					return
				}
				datastory.Stats.ReplyCount += 1
				err = uc.storyRepo.UpdateStory(ctx, datastory)
				if err != nil {
					resultChan <- &res.FailedStoryResponse{
						StoryID:      req.StoryId,
						UserID:       req.UserId,
						ErrorMessage: err.Error(),
					}
					return
				}
			case 1:
				newstoryID, err := gocql.ParseUUID(req.StoryId)
				if err != nil {
					resultChan <- &res.FailedStoryResponse{
						StoryID:      req.StoryId,
						UserID:       req.UserId,
						ErrorMessage: err.Error(),
					}
					return
				}
				newUserid, err := gocql.ParseUUID(req.UserId)
				if err != nil {
					resultChan <- &res.FailedStoryResponse{
						StoryID:      req.StoryId,
						UserID:       req.UserId,
						ErrorMessage: err.Error(),
					}
					return
				}
				entityStoryView := &entity.StoryView{
					StoryID:         newstoryID,
					ViewerID:        newUserid,
					ViewerName:      req.Name,
					ViewerAvatarURL: req.Avatar,
					ViewedAt:        req.ViewedAt,
					InteractionType: req.InteractionType,
					Content:         req.Content,
				}
				err = uc.storyview.CreateStoryView(ctx, entityStoryView)
				if err != nil {
					resultChan <- &res.FailedStoryResponse{
						StoryID:      req.StoryId,
						UserID:       req.UserId,
						ErrorMessage: err.Error(),
					}
					return
				}
			}
			if err != nil {
				resultChan <- &res.FailedStoryResponse{
					StoryID:      req.StoryId,
					UserID:       req.UserId,
					ErrorMessage: err.Error(),
				}
			} else {
				resultChan <- nil
			}
		})
		if err != nil {
			resultChan <- &res.FailedStoryResponse{
				StoryID:      req.StoryId,
				UserID:       req.UserId,
				ErrorMessage: err.Error(),
			}
			return nil, err
		}
	}
	var failedResponses []*res.FailedStoryResponse
	for i := 0; i < workercount; i++ {
		result := <-resultChan
		if result != nil {
			failedResponses = append(failedResponses, result)
			// Có thể log lỗi hoặc thực hiện hành động khắc phục nếu cần thiết
			// Ví dụ: log.Printf("Error processing reply for story %s and user %s: %s", result.StoryID, result.UserID, result.ErrorMessage)
		}
	}
	return failedResponses, nil
}
