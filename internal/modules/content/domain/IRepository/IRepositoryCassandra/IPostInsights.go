package IRepositoryCassandra

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type IPostInsights interface {
	CreatePostInsightInitPost(ctx context.Context, PostID string) error
	CreatePostInsightInitPostBulk(ctx context.Context, PostIDs []string) error
	UpdatePostInsightInteraction(ctx context.Context, PostID string, interactionType string, count int) error
	UpdatePostInsightInteractionBulk(ctx context.Context, reqs []*req.UpdatePostInsightsInteractionReq) error
	UpdatePostInsightLifeTime(ctx context.Context, PostID string, metricType string, value float64) error
	UpdatePostInsightLifeTimeBulk(ctx context.Context, reqs []*req.UpdatePostInsightsLifeTimeReq) error
	GetPostInsightByPostID(ctx context.Context, PostID string) (*entity.PostInsight, error)
	GetPostInsightsByPostIDs(ctx context.Context, PostIDs []string) ([]*entity.PostInsight, error)
	DeletePostInsightByPostID(ctx context.Context, PostID string) error
	DeletePostInsightsByPostIDBulk(ctx context.Context, PostIDs []string) error
}
