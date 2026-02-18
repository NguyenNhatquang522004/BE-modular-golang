package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type IBlockRepository interface {
	CreateBlockUser(blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	DeleteBlockUser(blockerUserID uuid.UUID, blockedUserID uuid.UUID) error
	UpdateBlockUser(blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error
	IsBlocked(blockerUserID uuid.UUID, blockedUserID uuid.UUID) (bool, error)
	GetBlockedUsers(blockerUserID uuid.UUID, blockedUserID uuid.UUID) (*entity.UserBlock, error)
	GetPaginationTypeBlock(BlockerUserID uuid.UUID, cursor string, limit int ,blocktype enum.Type_Block) (*dto.PaginationRes, error)
	GetBlockIndiscriminate(requesterID uuid.UUID, recipientID uuid.UUID ,) (*entity.UserBlock, error)
}
