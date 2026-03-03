package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
)

type ConsumerDeleteRelationTarget struct {
	events events.EventBus
}

func NewConsumerDeleteRelationTarget(events events.EventBus) *ConsumerDeleteRelationTarget {
	return &ConsumerDeleteRelationTarget{
		events: events,
	}
}

func (c *ConsumerDeleteRelationTarget) ConsumerDeleteRelationTarget(ctx context.Context) error {
	// Implement the logic for consuming delete relation target events
	err := c.events.Subscribe(ctx, constants.TopicDeleteCommunityRelationTarget.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		return nil
	})
	if err != nil {
		return err
	}
	// Log or handle successful subscription if needed
	// log.Printf("Subscribed to topic: %s", constants.TopicDeleteCommunityRelationTarget.String()
	return nil
}

func (c *ConsumerDeleteRelationTarget) ConsumerFailedDeleteRelationTarget(ctx context.Context) error {
	// Implement the logic for handling failed delete relation target events
	return nil
}
