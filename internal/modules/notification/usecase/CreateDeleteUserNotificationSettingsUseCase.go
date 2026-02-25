package usecase

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/identityEvent/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
)

type CreateDeleteUserNotificationSettingsUseCase struct {
	userNotificationRepo IRepositoryMongodb.IUserNotificationSettingsRepository
	eventBus             events.EventBus
	pool                 IRepositoryShare.IWorkerPool
}

func NewCreateDeleteUserNotificationSettingsUseCase(userNotificationRepo IRepositoryMongodb.IUserNotificationSettingsRepository,
	eventBus events.EventBus,
	pool IRepositoryShare.IWorkerPool) *CreateDeleteUserNotificationSettingsUseCase {
	return &CreateDeleteUserNotificationSettingsUseCase{
		userNotificationRepo: userNotificationRepo,
		eventBus:             eventBus,
		pool:                 pool,
	}
}
func (uc *CreateDeleteUserNotificationSettingsUseCase) Execute(ctx context.Context) error {
	err := uc.eventBus.Subscribe(ctx, string(constants.TopicCreateUserNotificationSettings), func(ctx context.Context, event events.IntegrationEvent) error {
		errs := uc.pool.Run(ctx, func() {
			payload, ok := event.Payload.(notificationEvent.CreateUserNotificationSettingsPayload)
			if !ok {
				// Xử lý lỗi nếu payload không đúng kiểu
				return
			}
			switch event.Type {
			case string(constants.Created):
				safectx := context.WithoutCancel(ctx)
				// Logic để tạo user notification settings dựa trên payload.UserID
				var req = &req.CreateUserNotificationSettingReq{}
				req.UserID = payload.UserID
				entity := mapper.ToEntityUserNotificationSetting(req)
				errs := uc.userNotificationRepo.CreateUserNotificationSettings(safectx, entity)
				if errs != nil {
					// Xử lý lỗi khi tạo user notification settings nếu cần thiết
					log.Printf("❌ Lỗi khi tạo user notification settings cho UserID %s: %v", payload.UserID, errs)
				}
			case string(constants.Deleted):
				safectx := context.WithoutCancel(ctx)
				// Logic để xóa user notification settings dựa trên payload.UserID
				errs := uc.userNotificationRepo.DeleteUserNotificationSettings(safectx, payload.UserID)
				if errs != nil {
					// Xử lý lỗi khi xóa user notification settings nếu cần thiết
					log.Printf("❌ Lỗi khi xóa user notification settings cho UserID %s: %v", payload.UserID, errs)
				}
			default:
				// Xử lý các loại sự kiện khác nếu cần thiết
			}
		})
		return errs
	})
	if err != nil {
		// Xử lý lỗi khi đăng ký sự kiện nếu cần thiết
		return err
	}
	return nil
}
