package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IFriendshipsRepository interface {
	CreateFriendship(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) error
	UpdateFriendshipStatus(ctx context.Context, data *entity.Friendships) error
	DeleteHardFriendship(ctx context.Context, friendshipID uuid.UUID) error
	DeleteSoftFriendship(ctx context.Context, friendshipID uuid.UUID) error
	PanigationStatusFriendship(ctx context.Context, userID uuid.UUID, cursor string, limit int, status enum.StatusFriendship) (*dto.PaginationRes, error)
	GetFriendshipTableByTableId(ctx context.Context, friendshipID uuid.UUID) (*entity.Friendships, error)
	GetFriendshipTableBybidirectional(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error)
	GetFriendshipIndiscriminate(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error)
}
