package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type ConsumerCommentStats struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	comment   IRepositoryMongoDB.ICommentRepository
}

func NewConsumerCommentStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, comment IRepositoryMongoDB.ICommentRepository) *ConsumerCommentStats {
	return &ConsumerCommentStats{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		comment:   comment,
	}
}
func (c *ConsumerCommentStats) ConsumerCommentStats(ctx context.Context) error {
	
	return nil
}
func (c *ConsumerCommentStats) ConsumerFailedCommentStats(ctx context.Context) error {
	return nil
}
