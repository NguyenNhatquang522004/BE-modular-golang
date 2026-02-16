package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type IFollowersRepository interface {
	CreateFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID) error
	DeleteSoftFollowUser(follower uuid.UUID) error
	DeleteHardFollowUser(follower uuid.UUID) error
	DeleteBatchSoftFollowUser(followerID uuid.UUID, followedUserID uuid.UUID) error
	DeleteBatchHardFollowUser(followerID uuid.UUID, followedUserID uuid.UUID) error
	UpdatateMuteFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID, isMuted bool) error
	PaginationFollowers(userID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
	PaginationFolloweds(FollowedUserID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error)
	GetFollowerTableByID(follower uuid.UUID) (*entity.Followers, error)
	GetFollowerTableBybidirectional(followerUserID uuid.UUID, followedUserID uuid.UUID) (*entity.Followers, error)
	GetFollowerIndiscriminate(requesterID uuid.UUID, recipientID uuid.UUID) (*[]entity.Followers, error)
}
