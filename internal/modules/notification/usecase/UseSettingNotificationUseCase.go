package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
)

type UseSettingNotificationUseCase struct {
	events events.EventBus
}

func NewUseSettingNotificationUseCase(events events.EventBus) *UseSettingNotificationUseCase {
	return &UseSettingNotificationUseCase{
		events: events,
	}
}

func (uc *UseSettingNotificationUseCase) Execute(ctx context.Context, req *req.UseSettingNotificationRequest) error {
	// Logic xử lý cài đặt notification
	// Ví dụ: Lắng nghe sự kiện thay đổi cài đặt notification và cập nhật vào database
	err := uc.events.Publish(ctx, constants.TopicUserNotificationSettings.String(), req.UserID, req.EventType.String(), req.NotificationChangePayload)
	if err != nil {
		return err
	}
	return nil
}
