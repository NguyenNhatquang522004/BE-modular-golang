package consumerContent

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepostitoryMongodb"
)

type ConsumerContent struct {
	// Thêm các trường cần thiết cho ConsumerContent tại đây
	pool            IRepositoryShare.IWorkerPool
	mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository
	eventBus        events.EventBus
}

func NewConsumerContent(pool IRepositoryShare.IWorkerPool, mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository, eventBus events.EventBus) *ConsumerContent {
	return &ConsumerContent{
		pool:            pool,
		mediaAssetsRepo: mediaAssetsRepo,
		eventBus:        eventBus,
	}
}

// payload *mediaInContent.CreateMediaAssetsPayload
// Implement các phương thức của IConsumerContent tại đây
func (c *ConsumerContent) ConsumerPublishPostCreateDeleteMediaAssets(ctx context.Context) error {
	// Logic để xử lý sự kiện tạo media assets cho một post
	err := c.eventBus.Subscribe(ctx, string(constants.TopicContentPostPublishMediaAssets), func(ctx context.Context, event events.IntegrationEvent) error {
		payload := event.Payload.(*mediaInContent.CreateMediaAssetsPayload) // Ép kiểu payload về đúng dạng
		switch event.Type {
		case constants.Created.String():
			entity, ok := mediaInContent.CreateMediaAssetsPayloadtoEntityMediaAssets(payload)
			if ok != nil {
				// Xử lý lỗi khi chuyển đổi payload sang entity nếu cần thiết
			}
			err := c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				// Giả sử CreateBulkMediaAssets là một phương thức để tạo nhiều media assets cùng lúc
				success, faildocs, err := c.mediaAssetsRepo.CreateBulkMediaAssets(safectx, entity)
				if err != nil {
					log.Printf("❌ Lỗi khi tạo media assets: %v", err)
					log.Printf("❌ Số lượng media assets tạo thành công: %d", success)
					for _, faildoc := range faildocs {
						log.Printf("❌ Media asset thất bại: %v", faildoc.ID)
						log.Printf("❌ Lỗi chi tiết: %v", faildoc.Reason)
					}
					// Có thể thêm logic để xử lý lỗi, ví dụ: retry, log chi tiết hơn, hoặc gửi một sự kiện lỗi khác
				}
				// Có thể thêm logic để xử lý lỗi, ví dụ: retry, log chi tiết hơn, hoặc gửi một sự kiện lỗi khác
			})
			if err != nil && err != context.Canceled {
				// Xử lý lỗi khi chạy công việc trong pool nếu cần thiết
			}
			log.Printf("✅ Đã xử lý xong media asset với MediaID: %s", payload.Items[0].MediaID) // Log thông tin media asset đã xử lý

			c.pool.Wait() // Đợi tất cả công việc trong pool hoàn thành trước khi tiếp tục

			return nil
		case constants.Deleted.String():
			err := c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				// Giả sử CreateBulkMediaAssets là một phương thức để tạo nhiều media assets cùng lúc
				err := c.mediaAssetsRepo.DeleteMediaAssetsByPostIDAndUserID(safectx, payload.PostID, payload.UserID)
				if err != nil {
					log.Printf("❌ Lỗi khi xóa media assets: %v", err)
				}
			})
			if err != nil && err != context.Canceled {

			}
			log.Printf("postid %s userid : %s", payload.PostID, payload.UserID) // Log thông tin media asset đã xử lý
			c.pool.Wait()
			return nil
		default:
			log.Printf("⚠️ Không xử lý được loại sự kiện: %s", event.Type)
			return nil
		}
	})
	if err != nil && err != context.Canceled {
		// Xử lý lỗi nếu cần thiết
		return err
	}
	return nil
}
