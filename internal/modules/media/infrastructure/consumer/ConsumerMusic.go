package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMusic struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	musicRepo IRepositoryMongodb.IMusicLibraryRepository
}

func NewConsumerMusic(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, musicRepo IRepositoryMongodb.IMusicLibraryRepository) *ConsumerMusic {
	return &ConsumerMusic{
		events:    events,
		pool:      pool,
		redisRepo: redisRepo,
		musicRepo: musicRepo,
	}
}

func (c *ConsumerMusic) ConsumerMusic(ctx context.Context) error {

	return nil
}

func (c *ConsumerMusic) ConsumerFailedMusic(ctx context.Context) error {
	return nil
}
