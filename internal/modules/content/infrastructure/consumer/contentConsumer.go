package content

import (
	"context"
	"errors"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/utils"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		data, ok := event.Payload.(*contentEvent.DeletePostPayload)
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
		data, ok := event.Payload.(*interactionEvent.EntityReactionPayload)
		if !ok {
			log.Printf("❌ Lỗi khi chuyển đổi payload: %v", event.Payload)
			return nil // Hoặc trả về lỗi nếu muốn dừng việc xử lý
		}
		switch event.Type {
		case constants.Created.String():
			// Xử lý khi có phản ứng mới
			datapost, err := c.postRepo.GetPostByID(ctx, data.TargetID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				utils.MappingCreateReactioncode(datapost, &data.ReactionCode)
				_, err = c.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(safectx, data.TargetID)
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
			datapost, err := c.postRepo.GetPostByID(ctx, data.TargetID)
			if err != nil {
				log.Printf("❌ Lỗi khi lấy bài viết: %v", err)
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				utils.MappingDeleteReactioncode(datapost, &data.ReactionCode)
				_, err := c.postRepo.UpdatePost(safectx, datapost)
				if err != nil {
					log.Printf("❌ Lỗi khi cập nhật bài viết: %v", err)
					return
				}
				datainsight, err := c.insightRepo.GetPostInsightByPostID(safectx, data.TargetID)
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
func (c *ContentConsumer) ConsumerSharePost(ctx context.Context) error {
	type taskResult struct {
		datapost *entity.Post
		data     *entity.PostExtension
		err      error
	}
	type taskResult2 struct {
		err error
	}
	worker := 2
	worker2 := 5
	resultChan := make(chan taskResult, worker)
	resultChan2 := make(chan taskResult2, worker2)
	err := c.eventBus.Subscribe(ctx, constants.TopicSharePost.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*contentEvent.SharePostPayload)
		if !ok {
			log.Printf("❌ Lỗi khi chuyển đổi payload: %v", event.Payload)
			return nil // Hoặc trả về lỗi nếu muốn dừng việc xử lý
		}
		for i := 0; i < worker; i++ {
			err := c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				switch i {
				case 0:
					datapost, err := c.postRepo.GetPostByID(safectx, data.PostID)
					resultChan <- taskResult{datapost: datapost, err: err}
				case 1:
					postExtension, err := c.extensionRepo.GetByPostID(safectx, data.PostID)
					resultChan <- taskResult{data: postExtension, err: err}
				}
			})
			if err != nil {
				return err
			}
		}
		c.pool.Wait() // Đảm bảo tất cả tác vụ đã hoàn thành trước khi tiếp tục xử lý kết quả
		var datapost *entity.Post
		var postExtension *entity.PostExtension
		for i := 0; i < worker; i++ {
			result := <-resultChan
			if result.err != nil {
				return result.err
			}
			if result.datapost != nil {
				datapost = result.datapost
			}
			if result.data != nil {
				postExtension = result.data
			}
		}
		if datapost == nil || postExtension == nil {
			return errors.New("failed to retrieve post data")
		}
		captureddatapost := datapost
		capturedpostExtension := postExtension
		newid := primitive.NewObjectID()
		for i := 0; i < worker2; i++ {
			err := c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				switch i {
				case 0:
					datapost.Stats.Shares += 1
					_, err := c.postRepo.UpdatePost(safectx, datapost)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
						return
					}
				case 1:
					datapost.UserID = data.UserID
					datapost.ID = newid
					_, err := c.postRepo.CreatePost(safectx, datapost)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
						return
					}
				case 2:
					postExtension.PostID = newid
					newidExtension := primitive.NewObjectID()
					postExtension.ID = newidExtension
					postExtension.ShareData.OriginalPostID = capturedpostExtension.ShareData.OriginalPostID
					postExtension.ShareData.ParentPostID = captureddatapost.ID
					postExtension.ShareData.Snapshot = entity.ShareSnapshot{
						AuthorID:       captureddatapost.UserID,
						AuthorName:     capturedpostExtension.ShareData.Snapshot.AuthorName,
						AuthorAvatar:   capturedpostExtension.ShareData.Snapshot.AuthorAvatar,
						ContentExcerpt: capturedpostExtension.ShareData.Snapshot.ContentExcerpt,
						MediaThumb:     capturedpostExtension.ShareData.Snapshot.MediaThumb,
						CreatedAt:      captureddatapost.CreatedAt,
					}

					err := c.extensionRepo.CreatePostExtension(safectx, postExtension)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
						return
					}
				case 3:
					idinsight := newid.String()
					postIDUUID, err := gocql.ParseUUID(idinsight)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
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
					err = c.insightRepo.CreatePostInsightInitPost(safectx, createPostInsight)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
						return
					}
				case 4:
					createSetting := &entity.PostSetting{
						ID:     primitive.NewObjectID(),
						PostID: newid,
					}
					_, err := c.settingRepo.CreatePostSetting(safectx, createSetting)
					if err != nil {
						resultChan2 <- taskResult2{err: err}
						return
					}
				}
			})
			if err != nil {
				resultChan2 <- taskResult2{err: err}
				return err
			}
		}
		c.pool.Wait() // Đảm bảo tất cả tác vụ đã hoàn thành trước khi tiếp tục xử lý kết quả
		for i := 0; i < worker2; i++ {
			result := <-resultChan2
			if result.err != nil {
				log.Printf("❌ Lỗi khi xử lý tác vụ chia sẻ bài viết: %v", result.err)
				return result.err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("❌ Lỗi Content Consumer: %v", err)
		return err
	}
	c.pool.Wait() // Đảm bảo tất cả tác vụ đã hoàn thành trước khi trả về
	return nil
}
