package consumer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroup struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	groupRepo IRepositoryMongodb.IGroupRepository
}

func NewConsumerGroup(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupRepo IRepositoryMongodb.IGroupRepository) *ConsumerGroup {
	return &ConsumerGroup{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		groupRepo: groupRepo,
	}
}
func (c *ConsumerGroup) ConsumerGroup(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGroup.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerGroup) ConsumerFailedGroup(ctx context.Context) error {

	return nil
}
