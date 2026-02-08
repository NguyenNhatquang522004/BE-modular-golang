package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/events"
)

type FeedConsumer struct {
	eventBus events.EventBus
}

func NewFeedConsumer(eb events.EventBus) *FeedConsumer {
	return &FeedConsumer{eventBus: eb}
}

// Start để đăng ký lắng nghe
func (c *FeedConsumer) Start(ctx context.Context) {
	err := c.eventBus.Subscribe(ctx, "user.registered", c.HandleUserRegistered)
	if err != nil {
		log.Fatalf("Failed to subscribe topic: %v", err)
	}
}

func (c *FeedConsumer) HandleUserRegistered(ctx context.Context, event events.IntegrationEvent) error {
	// Ép kiểu payload về đúng dạng
	payloadMap, ok := event.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload")
	}

	log.Printf("Creating welcome feed for User: %v", payloadMap["user_id"])
	// Gọi Logic domain của Feed module tại đây...
	return nil
}
