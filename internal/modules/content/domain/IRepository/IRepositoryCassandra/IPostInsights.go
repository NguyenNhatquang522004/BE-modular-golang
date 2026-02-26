package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostInsights interface {
	CreatePostInsightInitPost(ctx context.Context, postinsight *entity.PostInsight) error
	CreatePostInsightInitPostBulk(ctx context.Context, postinsights []*entity.PostInsight) (int64, []*cassandraErrors.InsightBulkError, error)
	UpdatePostInsightInteraction(ctx context.Context, PostID string, interactionType string, count int) error
	// UpdatePostInsightInteractionBulk(ctx context.Context, reqs []*req.UpdatePostInsightsInteractionReq) (int64, []*cassandraErrors.InsightBulkError, error)
	UpdatePostInsightLifeTime(ctx context.Context, PostID string, metricType string, value float64) error
	// UpdatePostInsightLifeTimeBulk(ctx context.Context, reqs []*req.UpdatePostInsightsLifeTimeReq) (int64, []*cassandraErrors.InsightBulkError, error)
	GetPostInsightByPostID(ctx context.Context, PostID string) (*entity.PostInsight, error)
	UpdateEntityPostInsight(ctx context.Context, postinsight *entity.PostInsight) error
	GetPostInsightsByPostIDs(ctx context.Context, PostIDs []string) ([]*entity.PostInsight, error)
	DeletePostInsightByPostID(ctx context.Context, PostID string) error
	DeletePostInsightsByPostIDBulk(ctx context.Context, PostIDs []string) (int64, []*cassandraErrors.InsightBulkError, error)
}
