package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
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
func (r *BlockRepository) CreateBlock(ctx context.Context, req *entity.UserBlock) error {
	return r.db.WithContext(ctx).Model(&entity.UserBlock{}).Create(req).Error
}
func (r *BlockRepository) UpdateBlock(ctx context.Context, req *entity.UserBlock) error {
	return r.db.WithContext(ctx).Model(&entity.UserBlock{}).Where(&entity.UserBlock{ID: req.ID}).Updates(req).Error
}
func (r *BlockRepository) DeleteBlock(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Model(&entity.UserBlock{}).Where(&entity.UserBlock{ID: uuid.MustParse(ID)}).Delete(&entity.UserBlock{}).Error
}
func (r *BlockRepository) GetBlockByID(ctx context.Context, ID string) (*entity.UserBlock, error) {
	var block *entity.UserBlock
	err := r.db.WithContext(ctx).Model(&entity.UserBlock{}).Where(&entity.UserBlock{ID: uuid.MustParse(ID)}).First(&block).Error
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (r *BlockRepository) GetPaginationTypeBlock(ctx context.Context, BlockerUserID uuid.UUID, cursor string, limit int, blocktype sharedEnums.Type_Block) (*dto.PaginationRes, error) {
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
