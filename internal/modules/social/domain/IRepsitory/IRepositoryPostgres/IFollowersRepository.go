package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type IFollowersRepository interface {
	CreateFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID) error
	DeleteSoftFollowUser(follower uuid.UUID) error
	DeleteHardFollowUser(follower uuid.UUID) error
	UpdatateMuteFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID, isMuted bool) error
	PaginationFollowers(userID uuid.UUID, cursor string, limit int) (*response.Response, error)
	PaginationFolloweds(FollowedUserID uuid.UUID, cursor string, limit int) (*response.Response, error)
	GetFollowerByID(follower uuid.UUID) (*entity.Followers, error)
	GetFollowerBybidirectional(followerUserID uuid.UUID, followedUserID uuid.UUID) (*entity.Followers, error)
}
