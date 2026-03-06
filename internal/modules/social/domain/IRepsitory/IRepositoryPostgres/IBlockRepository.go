package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IBlockRepository interface {
	CreateBlock(ctx context.Context, req *entity.UserBlock) error
	UpdateBlock(ctx context.Context, req *entity.UserBlock) error
	DeleteBlock(ctx context.Context, ID string) error
	GetBlockByID(ctx context.Context, ID string) (*entity.UserBlock, error)
	GetPaginationTypeBlock(ctx context.Context, BlockerUserID uuid.UUID, cursor string, limit int, blocktype enum.Type_Block) (*dto.PaginationRes, error)
}
