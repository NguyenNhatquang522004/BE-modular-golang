package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
)

type IFollowersRepository interface {
	CreateFollower(ctx context.Context, req *entity.Followers) error
	UpdateFollower(ctx context.Context, req *entity.Followers) error
	DeleteFollower(ctx context.Context, ID string) error
	DeleteFollowerByUserIDs(ctx context.Context, FollowerUserID string, FollowedUserID string) error
	GetFollowerByID(ctx context.Context, ID string) (*entity.Followers, error)
	GetFollowerByUserIDs(ctx context.Context, FollowerUserID string, FollowedUserID string) (*entity.Followers, error)
	PaginationFollowers(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	PaginationFolloweds(ctx context.Context, FollowedUserID string, cursor string, limit int) (*dto.PaginationRes, error)
}
