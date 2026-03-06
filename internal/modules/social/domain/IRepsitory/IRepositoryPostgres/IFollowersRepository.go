package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type IFollowersRepository interface {
	CreateFollower(ctx context.Context, req *entity.Followers) error
	UpdateFollower(ctx context.Context, req *entity.Followers) error
	DeleteFollower(ctx context.Context, ID string) error
	DeleteFollowerByUserID(ctx context.Context, FollowerUserID uuid.UUID, FollowedUserID uuid.UUID) error
	GetFollowerByID(ctx context.Context, ID string) (*entity.Followers, error)
	GetFollowerByUserIDs(ctx context.Context, FollowerUserID uuid.UUID, FollowedUserID uuid.UUID) (*entity.Followers, error)
	PaginationFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
	PaginationFolloweds(ctx context.Context, FollowedUserID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
}
