package content

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type ContentConsumer struct {
	eventBus      events.EventBus
	pool          IRepositoryShare.IWorkerPool
	postRepo      IRepositoryMongodb.IPostRepository
	mediaRepo     IRepositoryMongodb.IPostMediaRepository
	extensionRepo IRepositoryMongodb.IPostExtensionRepository
	settingRepo   IRepositoryMongodb.IPostSettingRepository
	insightRepo   IRepositoryCassandra.IPostInsights
}

// NewContentConsumer khởi tạo một instance của ContentConsumer
func NewContentConsumer(eventBus events.EventBus) *ContentConsumer {
	return &ContentConsumer{
		// Khởi tạo các trường cần thiết, ví dụ:
		// eventBus: eventBus,
	}
}

// Implement các phương thức của IContentConsumer tại đây
func (c *ContentConsumer) ConsumeContentPostStats(ctx context.Context) error {
	// Logic để tiêu thụ sự kiện liên quan đến thống kê bài viết
	return nil
}
func (c *ContentConsumer) ConsumeContentDeletePublishPost(ctx context.Context) error {
	err := c.eventBus.Subscribe(ctx, constants.TopicContentPostPublish.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*contentEvent.PostDeletePayload)
		if !ok {
			log.Printf("❌ Lỗi khi chuyển đổi payload: %v", event.Payload)
			return nil // Hoặc trả về lỗi nếu muốn dừng việc xử lý
		}
		log.Printf("✅ Nhận được sự kiện xóa bài viết với PostID: %s", data.PostID)
		err := c.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			c.postRepo.DeletePost(safectx, data.PostID)
			c.mediaRepo.DeleteByPostID(safectx, data.PostID)
			c.extensionRepo.DeleteByPostID(safectx, data.PostID)
			c.settingRepo.DeletePostSetting(safectx, data.PostID)
			c.insightRepo.DeletePostInsightByPostID(safectx, data.PostID)
		})
		return err
	})
	if err != nil && err != context.Canceled {
		log.Printf("❌ Lỗi Content Consumer: %v", err)
	}
	return nil
}
func (c *ContentConsumer) ConsumerReactionPost(ctx context.Context) error {
	err := c.eventBus.Subscribe(ctx, constants.TopicReactPost.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*interactionEvent.PostReactionPayload)
		if !ok {
			log.Printf("❌ Lỗi khi chuyển đổi payload: %v", event.Payload)
			return nil // Hoặc trả về lỗi nếu muốn dừng việc xử lý
		}
		switch event.Type {
		case constants.Created.String():
			// Xử lý khi có phản ứng mới
			datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				mappingCreateReactioncode(datapost, &data.ReactionCode)
				_, err = c.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(safectx, data.PostID)
				if err != nil {
					log.Printf("❌ Lỗi khi lấy insight bài viết: %v", err)
					return
				}
				datainsight.ReactionsTotal += 1
				err = c.insightRepo.UpdateEntityPostInsight(safectx, datainsight)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật insight bài viết: %v", err)
					return
				}
			})
			if err != nil {
				log.Printf("❌ Lỗi khi xử lý phản ứng mới: %v", err)
				return err
			}
		case constants.Deleted.String():
			// Xử lý khi có phản ứng bị xóa
			datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				mappingDeleteReactioncode(datapost, &data.ReactionCode)
				_, err := c.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(safectx, data.PostID)
				if err != nil {
					log.Printf("❌ Lỗi khi lấy insight bài viết: %v", err)
					return
				}
				datainsight.ReactionsTotal -= 1
				err = c.insightRepo.UpdateEntityPostInsight(safectx, datainsight)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật insight bài viết: %v", err)
					return
				}
			})
		}
		return nil
	})
	c.pool.Wait() // Đảm bảo tất cả tác vụ đã hoàn thành trước khi trả về
	if err != nil {
		log.Printf("❌ Lỗi Content Consumer: %v", err)
		return err
	}
	return nil
}
func (c *ContentConsumer) ConsumerCommentPost(ctx context.Context) error {
	err := c.eventBus.Subscribe(ctx, constants.TopicCounterPost.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*interactionEvent.CounterPostPayload)
		if !ok {
			log.Printf("❌ Lỗi khi chuyển đổi payload: %v", event.Payload)
			return nil // Hoặc trả về lỗi nếu muốn dừng việc xử lý
		}
		switch event.Type {
		case constants.Created.String():
			// Xử lý khi có bình luận mới
			datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				datapost.Stats.Comments += 1
				_, err = c.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(safectx, data.PostID)
				if err != nil {
					log.Printf("❌ Lỗi khi lấy insight bài viết: %v", err)
					return
				}
				datainsight.CommentsTotal += 1
				err = c.insightRepo.UpdateEntityPostInsight(safectx, datainsight)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật insight bài viết: %v", err)
					return
				}
			})

		case constants.Deleted.String():
			// Xử lý khi có bình luận bị xóa
			datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				datapost.Stats.Comments -= 1
				_, err = c.postRepo.UpdatePost(ctx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(ctx, data.PostID)
				if err != nil {
					log.Printf("❌ Lỗi khi lấy insight bài viết: %v", err)
					return
				}
				datainsight.CommentsTotal -= 1
				err = c.insightRepo.UpdateEntityPostInsight(ctx, datainsight)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật insight bài viết: %v", err)
					return
				}
			})
		}
		return nil
	})
	c.pool.Wait() // Đảm bảo tất cả tác vụ đã hoàn thành trước khi trả về
	if err != nil {
		log.Printf("❌ Lỗi Content Consumer: %v", err)
		return err
	}
	return nil
}

func mappingCreateReactioncode(datacomment *entity.Post, react *sharedEnums.ReactionCode) {
	datacomment.Stats.TotalReactions += 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Stats.Like += 1
	case sharedEnums.ReactionLove:
		datacomment.Stats.Love += 1
	case sharedEnums.ReactionHaha:
		datacomment.Stats.Haha += 1
	case sharedEnums.ReactionSad:
		datacomment.Stats.Sad += 1
	case sharedEnums.ReactionAngry:
		datacomment.Stats.Angry += 1
	case sharedEnums.ReactionWow:
		datacomment.Stats.Wow += 1
	}
}

func mappingDeleteReactioncode(datacomment *entity.Post, react *sharedEnums.ReactionCode) {
	datacomment.Stats.TotalReactions -= 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Stats.Like -= 1
	case sharedEnums.ReactionLove:
		datacomment.Stats.Love -= 1
	case sharedEnums.ReactionHaha:
		datacomment.Stats.Haha -= 1
	case sharedEnums.ReactionSad:
		datacomment.Stats.Sad -= 1
	case sharedEnums.ReactionAngry:
		datacomment.Stats.Angry -= 1
	case sharedEnums.ReactionWow:
		datacomment.Stats.Wow -= 1
	}
}
