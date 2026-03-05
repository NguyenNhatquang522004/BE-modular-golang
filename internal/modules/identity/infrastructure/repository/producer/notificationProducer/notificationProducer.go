package notificationProducer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/identityEvent/notificationEvent"
)

type NotificationProducer struct {
	// Có thể thêm các trường như Kafka producer, RabbitMQ producer, v.v. tùy vào hệ thống message queue bạn sử dụng
	eventBus events.EventBus
}

func NewNotificationProducer(eventBus events.EventBus) *NotificationProducer {
	return &NotificationProducer{
		eventBus: eventBus,
	}
}

func (p *NotificationProducer) CreateUserNotificationSettings(ctx context.Context, payload *notificationEvent.CreateUserNotificationSettingsPayload) error {
	// Thực hiện logic gửi thông báo ở đây, ví dụ: gửi message đến Kafka, RabbitMQ, v.v.
	err := p.eventBus.Publish(ctx, constants.TopicUserNotificationSettings.String(), payload.UserID, constants.Created.String(), payload)
	if err != nil {
		return err
	}
	return nil
}
