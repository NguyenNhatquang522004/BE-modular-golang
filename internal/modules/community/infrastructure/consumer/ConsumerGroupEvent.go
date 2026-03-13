package consumer

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerGroupEvent struct {
	events         events.EventBus
	pool           IRepositoryShare.IWorkerPool
	redisRepo      IRepositoryShare.IRedis
	groupEventRepo IRepositoryMongodb.IGroupEventsRepository
}

func NewConsumerGroupEvent() *ConsumerGroupEvent {
	return &ConsumerGroupEvent{}
}
