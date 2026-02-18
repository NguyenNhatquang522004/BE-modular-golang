package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IFriendshipsRepository interface {
	CreateFriendship(Requester_ID uuid.UUID, Recipient_ID uuid.UUID) error
	UpdateFriendshipStatus(data *entity.Friendships) error
	DeleteHardFriendship(friendshipID uuid.UUID) error
	DeleteSoftFriendship(friendshipID uuid.UUID) error
	PanigationStatusFriendship(userID uuid.UUID, cursor string, limit int, status enum.StatusFriendship) (*dto.PaginationRes, error)
	GetFriendshipTableByTableId(friendshipID uuid.UUID) (*entity.Friendships, error)
	GetFriendshipTableBybidirectional(requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error)
	GetFriendshipIndiscriminate(requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error)
}
