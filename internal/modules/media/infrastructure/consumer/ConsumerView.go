package consumer

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/gocql/gocql"
)

type ConsumerView struct {
	events        events.EventBus
	storyViewRepo IRepositoryCassandra.IStoryViewRepository
	storyRepo     IRepositoryMongodb.IStoryRepository
	pool          IRepositoryShare.IWorkerPool
}

func NewConsumerView(events events.EventBus, storyViewRepo IRepositoryCassandra.IStoryViewRepository, storyRepo IRepositoryMongodb.IStoryRepository, pool IRepositoryShare.IWorkerPool) *ConsumerView {
	return &ConsumerView{
		events:        events,
		storyViewRepo: storyViewRepo,
		storyRepo:     storyRepo,
		pool:          pool,
	}
}
func (c *ConsumerView) ConsumerViewCountStory(ctx context.Context) {
	type taskResult struct {
		storyId string
		userId  string
		err     error
	}
	worker := 2
	err := c.events.Subscribe(ctx, constants.TopicViewCountStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*req.ViewCountStoryRequest)
		dataStory, err := c.storyRepo.GetStoryByID(ctx, data.StoryId)
		if err != nil {
			return err
		}
		resultChan := make(chan taskResult, worker)
		for i := 0; i < worker; i++ {
			err = c.pool.Run(ctx, func() {
				switch i {
				case 0:
					dataStory.Stats.ViewsCount += 1
					if len(dataStory.PreviewViewers) == 3 {
						dataStory.PreviewViewers = dataStory.PreviewViewers[1:]
						dataStory.PreviewViewers = append(dataStory.PreviewViewers, entity.ViewerPreview{
							UserID: data.UserId,
							Avatar: data.Avatar,
							Name:   data.Name,
						})
					} else {
						dataStory.PreviewViewers = append(dataStory.PreviewViewers, entity.ViewerPreview{
							UserID: data.UserId,
							Avatar: data.Avatar,
							Name:   data.Name,
						})
					}
					err = c.storyRepo.UpdateStory(ctx, dataStory)
					if err != nil {
						// Xử lý lỗi nếu cần thiết
						resultChan <- taskResult{
							storyId: data.StoryId,
							userId:  data.UserId,
							err:     err,
						}
						return
					}
				case 1:
					gocqlid, err := gocql.ParseUUID(data.StoryId)
					if err != nil {
						// Xử lý lỗi nếu cần thiết
						resultChan <- taskResult{
							storyId: data.StoryId,
							userId:  data.UserId,
							err:     err,
						}
						return
					}
					cqluserid, err := gocql.ParseUUID(data.UserId)
					if err != nil {
						// Xử lý lỗi nếu cần thiết
						resultChan <- taskResult{
							storyId: data.StoryId,
							userId:  data.UserId,
							err:     err,
						}
						return
					}
					err = c.storyViewRepo.CreateStoryView(ctx, &entity.StoryView{
						StoryID:         gocqlid,
						ViewerID:        cqluserid,
						ViewerAvatarURL: data.Avatar,
						ViewerName:      data.Name,
						ViewedAt:        data.ViewedAt,
						InteractionType: data.InteractionType,
						ReactionCode:    data.ReactionCode,
						PollOptionIndex: data.PollOptionIndex,
						Content:         data.Content,
					})
					if err != nil {
						// Xử lý lỗi nếu cần thiết
						resultChan <- taskResult{
							storyId: data.StoryId,
							userId:  data.UserId,
							err:     err,
						}
						return
					}
					// Có thể thêm logic cập nhật cache hoặc các hệ thống khác nếu cần thiết
					// Ví dụ: Cập nhật cache tổng số lượt xem của story
				}
			})
		}
		for i := 0; i < worker; i++ {
			result := <-resultChan
			if result.err != nil {
				err := c.events.Publish(ctx, constants.TopicViewCountStory.String(), result.storyId, constants.Deleted.String(), &res.FailedStoryResponse{
					StoryID:      result.storyId,
					UserID:       result.userId,
					ErrorMessage: result.err.Error(),
				})
				if err != nil {
					// Xử lý lỗi nếu cần thiết
					log.Printf("Error publishing failed view count event for story %s and user %s: %v", result.storyId, result.userId, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		// Xử lý lỗi nếu cần thiết
		log.Printf("Error subscribing to view count story events: %v", err)
	}
	return

}
func (c *ConsumerView) ConsumerFailedViewCountStory(ctx context.Context) {
	err := c.events.Subscribe(ctx, constants.TopicViewCountStory.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*res.FailedStoryResponse)
		// Xử lý logic khi nhận được sự kiện thất bại, ví dụ: ghi log, gửi thông báo, v.v.
		log.Printf("Failed to update view count for story %s and user %s: %s", data.StoryID, data.UserID, data.ErrorMessage)
		return nil
	})
	if err != nil {
		// Xử lý lỗi nếu cần thiết
		log.Printf("Error subscribing to failed view count story events: %v", err)
	}
}
