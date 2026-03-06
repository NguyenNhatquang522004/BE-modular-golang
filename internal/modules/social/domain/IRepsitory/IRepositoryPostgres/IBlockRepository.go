package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
)

type IBlockRepository interface {
	CreateBlock(ctx context.Context, req *entity.UserBlock) error
	UpdateBlock(ctx context.Context, req *entity.UserBlock) error
	DeleteBlockByUserIDs(ctx context.Context, blockerUserID string, blockedUserID string) error
	DeleteBlock(ctx context.Context, ID string) error
	GetListBlockByUserID(ctx context.Context, userID string) ([]*entity.UserBlock, error)
	GetBlockByUserIDs(ctx context.Context, blockerUserID string, blockedUserID string) (*entity.UserBlock, error)
	GetBlockByID(ctx context.Context, ID string) (*entity.UserBlock, error)
	GetPaginationTypeBlock(ctx context.Context, BlockerUserID uuid.UUID, cursor string, limit int, blocktype sharedEnums.Type_Block) (*dto.PaginationRes, error)
}
