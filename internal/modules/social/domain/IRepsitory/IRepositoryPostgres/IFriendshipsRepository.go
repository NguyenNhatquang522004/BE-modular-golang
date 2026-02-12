package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IFriendshipsRepository interface {
	CreateFriendship(Requester_ID uuid.UUID, Recipient_ID uuid.UUID) error
	UpdateFriendshipStatus(friendshipID uuid.UUID, status enum.StatusFriendship) error
	DeleteHardFriendship(friendshipID uuid.UUID) error
	DeleteSoftFriendship(friendshipID uuid.UUID) error
	PanigationAcceptedFriendship(userID uuid.UUID, cursor string, limit int) (*response.Response, error)
	PanigationPendingFriendship(userID uuid.UUID, cursor string, limit int) (*response.Response, error)
	GetFriendshipByUserIDs(friendshipID uuid.UUID) (*entity.Friendships, error)
	GetFriendshipBybidirectional(requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error)
}
