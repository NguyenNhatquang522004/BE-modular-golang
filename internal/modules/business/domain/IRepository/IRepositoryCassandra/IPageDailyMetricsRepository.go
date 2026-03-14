package IRepositoryCassandra

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IPageDailyMetricsRepository interface {
	// Define the methods for the PageDailyMetricsRepository interface here
	CreatePageDailyMetric(ctx context.Context, metric *entity.PageDailyMetric) error
	CreateBulkPageDailyMetrics(ctx context.Context, metrics []*entity.PageDailyMetric) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error)
	GetPageDailyMetricsByPageID(ctx context.Context, pageID string, metricDate time.Time, cursor string, limit int) (*dto.PaginationRes, error)
	UpdatePageDailyMetric(ctx context.Context, metric *entity.PageDailyMetric) error
	UpdateBulkPageDailyMetrics(ctx context.Context, metrics []*entity.PageDailyMetric) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error)
	DeletePageDailyMetricsByPageID(ctx context.Context, pageID string) error
	DeletePageDailyMetric(ctx context.Context, pageID string, metricDate time.Time) error
	DeleteBulkPageDailyMetrics(ctx context.Context, pageID string, metricDates []time.Time) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error)
	InsertOrUpdatePageDailyMetric(ctx context.Context, metric *entity.PageDailyMetric) error
	GetPageDailyMetricLatest(ctx context.Context, pageID string) (*entity.PageDailyMetric, error)
}
