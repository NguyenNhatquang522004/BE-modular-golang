package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type IFollowersRepository interface {
	CreateFollowUser(ctx context.Context, followerUserID uuid.UUID, followedUserID uuid.UUID) error
	DeleteSoftFollowUser(ctx context.Context, follower uuid.UUID) error
	DeleteHardFollowUser(ctx context.Context, follower uuid.UUID) error
	DeleteBatchSoftFollowUser(ctx context.Context, followerID uuid.UUID, followedUserID uuid.UUID) error
	DeleteBatchHardFollowUser(ctx context.Context, followerID uuid.UUID, followedUserID uuid.UUID) error
	UpdatateMuteFollowUser(ctx context.Context, followerUserID uuid.UUID, followedUserID uuid.UUID, isMuted bool) error
	PaginationFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
	PaginationFolloweds(ctx context.Context, FollowedUserID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
	GetFollowerTableByID(ctx context.Context, follower uuid.UUID) (*entity.Followers, error)
	GetFollowerTableBybidirectional(ctx context.Context, followerUserID uuid.UUID, followedUserID uuid.UUID) (*entity.Followers, error)
	GetFollowerIndiscriminate(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) ([]*entity.Followers, error)
}
