package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type ConsumerPageDailyMetrics struct {
	// Define the fields for the ConsumerPageDailyMetrics struct here
	events           events.EventBus
	pageRepo         IRepositoryMongodb.IPagesRepository
	pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository
	pageMetricRepo   IRepositoryCassandra.IPageDailyMetricsRepository
	pool             IRepositoryShare.IWorkerPool
}

func NewConsumerPageDailyMetrics(events events.EventBus, pageRepo IRepositoryMongodb.IPagesRepository, pageFollowerRepo IRepositoryMongodb.IPageFollowersRepository, pageMetricRepo IRepositoryCassandra.IPageDailyMetricsRepository, pool IRepositoryShare.IWorkerPool) *ConsumerPageDailyMetrics {
	return &ConsumerPageDailyMetrics{
		events:           events,
		pageRepo:         pageRepo,
		pageFollowerRepo: pageFollowerRepo,
		pageMetricRepo:   pageMetricRepo,
		pool:             pool,
	}
}

func (c *ConsumerPageDailyMetrics) ConsumerPageDailyMetric(ctx context.Context) error {
	// Implement the logic for consuming page daily metrics here
	return nil
}

func (c *ConsumerPageDailyMetrics) ConsumerFailedPageDailyMetric(ctx context.Context) error {
	// Implement the logic for consuming failed page daily metrics here
	return nil
}
