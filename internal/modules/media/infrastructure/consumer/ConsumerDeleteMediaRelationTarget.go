package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
)

type ConsumerDeleteMediaRelationTarget struct {
	events events.EventBus
}

func NewConsumerDeleteMediaRelationTarget(events events.EventBus) *ConsumerDeleteMediaRelationTarget {
	return &ConsumerDeleteMediaRelationTarget{
		events: events,
	}
}

func (c *ConsumerDeleteMediaRelationTarget) ConsumerDeleteMediaRelationTarget(ctx context.Context) error {

	err := c.events.Subscribe(ctx, constants.TopicDeleteMediaRelationTarget.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerDeleteMediaRelationTarget) ConsumerFailedDeleteMediaRelationTarget(ctx context.Context) error {
	return nil
}
