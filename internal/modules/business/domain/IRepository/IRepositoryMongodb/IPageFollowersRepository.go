package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IPageFollowersRepository interface {
	// Define the methods for the PageFollowersRepository interface here
	CreateFollower(ctx context.Context, follower *entity.PageFollower) error
	CreateBulkFollowers(ctx context.Context, followers []*entity.PageFollower) (int64, []*mongodbErrors.BulkError, error)
	GetFollowersByPageID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetFollowerByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateFollower(ctx context.Context, follower *entity.PageFollower) error
	UpdateBulkFollowers(ctx context.Context, followers []*entity.PageFollower) (int64, []*mongodbErrors.BulkError, error)
	DeleteFollower(ctx context.Context, pageID string, userID string) error
	DeleteBulkFollowers(ctx context.Context, pageID string, userIDs []string) (int64, []*mongodbErrors.BulkError, error)
	DeleteBulkFollowersByPageID(ctx context.Context, pageID string) (int64, []*mongodbErrors.BulkError, error)
	DeleteFollowerByPageID(ctx context.Context, pageID string) error
}
