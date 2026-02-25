package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IStrategy"
)

type SendNotificationTypeUseCase struct {
	eventBus        events.EventBus
	pool            IRepositoryShare.IWorkerPool
	handlerStrategy map[string]IStrategy.IStrategyTypeNotificationType
}

func NewSendNotificationTypeUseCase(eventBus events.EventBus, StrategyNotificationType []IStrategy.IStrategyTypeNotificationType, pool IRepositoryShare.IWorkerPool) *SendNotificationTypeUseCase {
	hmap := make(map[string]IStrategy.IStrategyTypeNotificationType)
	for _, strategy := range StrategyNotificationType {
		hmap[strategy.GetType().String()] = strategy
	}

	return &SendNotificationTypeUseCase{
		eventBus:        eventBus,
		handlerStrategy: hmap,
		pool:            pool,
	}
}

func (uc *SendNotificationTypeUseCase) Execute(ctx context.Context , ) error {
	// Example: Determine the notification type and execute the corresponding strategy
	err := uc.eventBus.Subscribe(ctx, string(constants.TopicSendNotificationType), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*notificationEvent.NotificationPayload)
		for _, typeNotification := range data.TypeNotification {
			strategy, exists := uc.handlerStrategy[typeNotification.String()]
			if !exists {
				// Handle the case where the strategy does not exist for the given notification type
				continue
			}
			err := uc.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				err := strategy.Execute(safectx, data)
				if err != nil {
					// Handle the error appropriately, e.g., log it or return an error
					return
				}
			})
			if err != nil {
				// Handle the error appropriately, e.g., log it or return an error
				return err
			}
		}
		uc.pool.Wait() // Wait for all goroutines to finish
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
