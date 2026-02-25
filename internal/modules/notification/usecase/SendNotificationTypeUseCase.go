package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IStrategy"
)

type SendNotificationTypeUseCase struct {
	handlerStrategy map[string]IStrategy.IStrategyTypeNotificationType
}

func NewSendNotificationTypeUseCase(StrategyNotificationType []IStrategy.IStrategyTypeNotificationType) *SendNotificationTypeUseCase {
	hmap := make(map[string]IStrategy.IStrategyTypeNotificationType)
	for _, strategy := range StrategyNotificationType {
		hmap[strategy.GetType().String()] = strategy
	}

	return &SendNotificationTypeUseCase{
		handlerStrategy: hmap,
	}
}

func (uc *SendNotificationTypeUseCase) Execute(ctx context.Context) (*response.Response, error) {
	// Example: Determine the notification type and execute the corresponding strategy
	return &response.Response{
		Message: "Notification sent successfully",
	}, nil
}
