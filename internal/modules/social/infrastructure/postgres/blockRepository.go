package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlockRepository struct {
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

// Implement the methods defined in IBlockRepository interface
func NewBlockRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *BlockRepository {
	return &BlockRepository{
		db:        db,
		redisRepo: redisRepo,
	}
}

func (r *BlockRepository) CreateBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error {
	// Implementation here

	return r.db.Create(&entity.UserBlock{
		Blocker_UserID: blockerUserID,
		Blocked_UserID: blockedUserID,
		Type_Block:     statusBlock,
	}).Error
}
func (r *BlockRepository) DeleteBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) error {
	// Implementation here
	return r.db.Where(&entity.UserBlock{Blocker_UserID: blockerUserID, Blocked_UserID: blockedUserID}).Delete(&entity.UserBlock{}).Error
}
func (r *BlockRepository) UpdateBlockUser(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID, statusBlock enum.Type_Block) error {
	// Implementation here
	data, err := r.GetBlockedUsers(ctx, blockerUserID, blockedUserID)
	if err != nil {
		return err
	}
	data.Type_Block = statusBlock
	err = r.db.Save(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BlockRepository) IsBlocked(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) (bool, error) {
	// Implementation here
	data, err := r.GetBlockIndiscriminate(ctx, blockerUserID, blockedUserID)
	if err != nil {
		return false, err
	}
	if data != nil {
		return true, nil
	}
	return false, nil
}

func (r *BlockRepository) GetBlockedUsers(ctx context.Context, blockerUserID uuid.UUID, blockedUserID uuid.UUID) (*entity.UserBlock, error) {
	// Implementation here
	var data = &entity.UserBlock{}
	err := r.db.Where(&entity.UserBlock{Blocker_UserID: blockerUserID, Blocked_UserID: blockedUserID}).Or(&entity.UserBlock{Blocker_UserID: blockedUserID, Blocked_UserID: blockerUserID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *BlockRepository) GetPaginationTypeBlock(ctx context.Context, BlockerUserID uuid.UUID, cursor string, limit int, blocktype enum.Type_Block) (*dto.PaginationRes, error) {
	// Implementation here
	data := []*entity.UserBlock{}
	querylimit := limit + 1
	itemset := []string{
		"block_cache_user_" + BlockerUserID.String(),
		"block_cache_hasnext_user_" + BlockerUserID.String(),
		"block_cache_nextcursor_user_" + BlockerUserID.String(),
		"block_cache_limit_user_" + BlockerUserID.String(),
	}
	cacheData, cacheCursor, hasNextCache, cacheLimit, err := r.redisRepo.CustomizeGetCache(ctx, itemset)
	if err != nil {
		return nil, err
	}
	if cacheData != nil {
		return &dto.PaginationRes{
			NextCursor: cacheCursor,
			HasNext:    hasNextCache,
			Data:       cacheData,
			Limit:      cacheLimit,
		}, nil
	}
	query := r.db.Where(&entity.UserBlock{Blocker_UserID: BlockerUserID, Type_Block: blocktype}).Order("created_at DESC ,id DESC").Limit(querylimit)
	if cursor != "" {
		time, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where(" (created_at < ?) OR (created_at = ? AND id > ?)", time, time, id)
		err = query.Find(data).Error
		if err != nil {
			return nil, err
		}
		var hasNext = false
		if len(data) > limit {
			hasNext = true
			data = data[:limit]
		}

		lastBlock := data[len(data)-1]
		nextCursor := utils.EncodeCursor(lastBlock.CreatedAt, lastBlock.ID)
		return &dto.PaginationRes{
			NextCursor: nextCursor,
			HasNext:    hasNext,
			Data:       data,
			Limit:      limit,
		}, nil
	}
	err = query.Find(data).Error
	if err != nil {
		return nil, err
	}
	var hasNext = false
	if len(data) > limit {
		hasNext = true
		data = data[:limit]
	}

	var nextCursor string
	if len(data) > 0 {
		lastBlock := data[len(data)-1]
		nextCursor = utils.EncodeCursor(lastBlock.CreatedAt, lastBlock.ID)
	}

	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       data,
		Limit:      limit,
	}, nil
}

func (r *BlockRepository) GetBlockIndiscriminate(ctx context.Context, requesterID uuid.UUID, recipientID uuid.UUID) (*entity.UserBlock, error) {
	data := &entity.UserBlock{}
	err := r.db.Where(&entity.UserBlock{Blocker_UserID: requesterID, Blocked_UserID: recipientID}).Or(&entity.UserBlock{Blocker_UserID: recipientID, Blocked_UserID: requesterID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
