package IRepositoryCassandra

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IStoryViewRepository interface {
	// Define methods for StoryViewRepository here
	CreateStoryView(ctx context.Context, storyView *entity.StoryView) error
	CreateBulkStoryViews(ctx context.Context, storyViews []*entity.StoryView) (int64, []*cassandraErrors.StoryViewBulkError, error)
	GetStoryViewsByStoryID(ctx context.Context, storyID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetStoryViewsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateStoryView(ctx context.Context, storyView *entity.StoryView) error
	UpdateBulkStoryViews(ctx context.Context, storyViews []*entity.StoryView) (int64, []*cassandraErrors.StoryViewBulkError, error)
	DeleteStoryView(ctx context.Context, storyID string) error
	DeleteBulkStoryViews(ctx context.Context, storyIDs []string) (int64, []*cassandraErrors.StoryViewBulkError, error)
	DeleteStoryViewsByUserID(ctx context.Context, userID string) error
	DeleteStoryViewsByUserIAndStoryID(ctx context.Context, userID string, storyID string) error
	DeleteBulkStoryViewsByUserIAndStoryID(ctx context.Context, userID string, storyID []string) (int64, []*cassandraErrors.StoryViewBulkError, error)
	DeleteStoryViewsByStoryIDAndUserID(ctx context.Context, storyID string, userID string) error
	DeleteBulkStoryViewsByStoryIDAndUserID(ctx context.Context, storyID []string, userID string) (int64, []*cassandraErrors.StoryViewBulkError, error)
	DeleteStoryViewsByViewedAt(ctx context.Context, storyID string, viewedAt time.Time) error
	DeleteBulkStoryViewsByViewedAt(ctx context.Context, storyIDs []string, viewedAt time.Time) (int64, []*cassandraErrors.StoryViewBulkError, error)
}
