package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
)

type ConsumerCommentLive struct {
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	livecomment IRepositoryCassandra.ILiveCommentsRepository
}

func NewConsumerCommentLive(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, livecomment IRepositoryCassandra.ILiveCommentsRepository) *ConsumerCommentLive {
	return &ConsumerCommentLive{
		events:      events,
		pool:        pool,
		redisRepo:   redisRepo,
		livecomment: livecomment,
	}
}
func (c *ConsumerCommentLive) ConsumerStartLiveStream(ctx context.Context) error {
	return nil
}
func (c *ConsumerCommentLive) ConsumerFailedStartLiveStream(ctx context.Context) error {
	return nil
}
