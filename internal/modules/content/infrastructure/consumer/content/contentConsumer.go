package content

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
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
	err := c.eventBus.Subscribe(ctx, constants.TopicContent.String(), func(ctx context.Context, event events.IntegrationEvent) error {
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
