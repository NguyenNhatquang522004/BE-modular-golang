package IRepostitoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IReelRepository interface {
	// Define methods for ReelRepository here
	CreateReel(ctx context.Context, reel *entity.Reel) error
	CreateBulkReels(ctx context.Context, reels []*entity.Reel) (int64, []*mongodbErrors.BulkError, error)
	GetReelByID(ctx context.Context, id string) (*entity.Reel, error)
	GetReelsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateReel(ctx context.Context, reel *entity.Reel) error
	UpdateBulkReels(ctx context.Context, reels []*entity.Reel) (int64, []*mongodbErrors.BulkError, error)
	DeleteReel(ctx context.Context, id string) error
	DeleteBulkReels(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
}
