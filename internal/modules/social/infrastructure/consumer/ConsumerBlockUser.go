package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
)

type ConsumerBlockUser struct {
	redisRepo IRepositoryShare.IRedis
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
}

func NewConsumerBlockUser() *ConsumerBlockUser {
	return &ConsumerBlockUser{}
}

func (c *ConsumerBlockUser) ConsumerBlockUser(ctx context.Context) error {
	// Implement the logic for consuming block user events
	return nil
}

func (c *ConsumerBlockUser) ConsumerFailedBlockUser(ctx context.Context) error {
	// Implement the logic for consuming failed block user events
	return nil
}
