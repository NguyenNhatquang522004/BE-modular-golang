package IStrategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type IStrategyTypeNotificationType interface {
	// Define method signatures for notification strategies
	Execute(ctx context.Context, req *notificationEvent.NotificationPayload) error
	GetType() sharedEnums.NotificationType
}
