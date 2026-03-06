package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IFriendshipsRepository interface {
	CreateFriendship(ctx context.Context, req *entity.Friendships) error
	UpdateFriendship(ctx context.Context, req *entity.Friendships) error
	DeleteFriendship(ctx context.Context, ID string) error
	GetFriendshipByID(ctx context.Context, ID string) (*entity.Friendships, error)
	PanigationStatusFriendship(ctx context.Context, userID uuid.UUID, cursor string, limit int, status enum.StatusFriendship) (*dto.PaginationRes, error)
}
