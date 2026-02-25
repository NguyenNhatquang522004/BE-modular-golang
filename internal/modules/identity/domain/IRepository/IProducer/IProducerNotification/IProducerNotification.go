package IProducerNotification

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/identityEvent/notificationEvent"
)

type IProducerNotification interface {
	CreateUserNotificationSettings(ctx context.Context, payload *notificationEvent.CreateUserNotificationSettingsPayload) error
}
