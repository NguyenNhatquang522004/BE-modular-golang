package IRepostitoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IStoryRepository interface {
	// Define methods for StoryRepository here
	CreateStory(ctx context.Context, story *entity.Story)  error
	CreateBulkStories(ctx context.Context, stories []*entity.Story) (int64, []*dto.BulkError, error)
	GetStoryByID(ctx context.Context, id string) (*entity.Story, error)
	GetBulkStoriesByID(ctx context.Context, ids []string) ([]*entity.Story, error)
	GetStoriesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateStory(ctx context.Context, story *entity.Story)  error
	UpdateBulkStories(ctx context.Context, stories []*entity.Story) (int64, []*dto.BulkError, error)
	DeleteStory(ctx context.Context, id string) error
	DeleteBulkStories(ctx context.Context, ids []string) (int64, []*dto.BulkError, error)
}
