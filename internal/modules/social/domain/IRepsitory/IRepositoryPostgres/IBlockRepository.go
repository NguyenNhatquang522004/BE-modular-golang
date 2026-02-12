package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IBlockRepository interface {
	CreateBlockUser(blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	UpdateBlockUser(blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	IsBlocked(blockerUserID uuid.UUID, blockedUserID uuid.UUID) (bool, error)
	GetBlockedUsers(blockerUserID uuid.UUID, blockedUserID uuid.UUID) ([]*entity.UserBlock, error)
	GetPaginationTypeBlock(BlockerUserID uuid.UUID, cursor string, limit int) (*response.Response, error)
}
