package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IBlockRepository interface {
	CreateBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	DeleteBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error
	UpdateBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	IsBlocked(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) (bool, error)
	GetBlockedUsers(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) (*entity.UserBlock, error)
	GetPaginationTypeBlock(ctx context.Context, BlockerUserID uuid.UUID, cursor string, limit int, blocktype enum.Type_Block) (*dto.PaginationRes, error)
	GetBlockIndiscriminate(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) (*entity.UserBlock, error)
}
