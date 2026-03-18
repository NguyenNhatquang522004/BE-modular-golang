package consumer

import (
	"context"
	"sync"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/IRepository/neo4j"
)

type ConsumerGraphFriendshipFrequencies struct {
	graphRepo neo4j.IGraphRepository
	pool      IRepositoryShare.IWorkerPool
	sheduler  IRepositoryShare.IScheduler
}

func NewConsumerGraphFriendshipFrequencies(graphRepo neo4j.IGraphRepository, pool IRepositoryShare.IWorkerPool, scheduler IRepositoryShare.IScheduler) *ConsumerGraphFriendshipFrequencies {
	return &ConsumerGraphFriendshipFrequencies{
		graphRepo: graphRepo,
		pool:      pool,
		sheduler:  scheduler,
	}
}

func (c *ConsumerGraphFriendshipFrequencies) ConsumerGraphFriendshipFrequencies(ctx context.Context) error {
	// Implement the logic for consuming graph friendship frequencies
	err := c.sheduler.ScheduleJob("0 */12 * * *", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		c.pool.Run(ctx, func() {
			defer wg.Done()
			err := c.graphRepo.SyncAllFriendshipFrequencies(ctx, 1000)
			if err != nil {
				// Log the error
			}
		})
		wg.Wait()
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerGraphFriendshipFrequencies) ConsumerFailedGraphFriendshipFrequencies(ctx context.Context) error {
	// Implement the logic for handling failed graph friendship frequencies
	return nil
}
