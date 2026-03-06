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

type FriendshipsRepository struct {
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

// Implement the methods defined in IFriendshipsRepository interface
func NewFriendshipsRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *FriendshipsRepository {
	return &FriendshipsRepository{db: db, redisRepo: redisRepo}
}
func (r *FriendshipsRepository) CreateFriendship(ctx context.Context, req *entity.Friendships) error {
	return r.db.WithContext(ctx).Model(&entity.Friendships{}).Create(req).Error
}
func (r *FriendshipsRepository) UpdateFriendship(ctx context.Context, req *entity.Friendships) error {
	return r.db.WithContext(ctx).Model(&entity.Friendships{}).Where(&entity.Friendships{ID: req.ID}).Updates(req).Error
}
func (r *FriendshipsRepository) DeleteFriendship(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Model(&entity.Friendships{}).Where(&entity.Friendships{ID: uuid.MustParse(ID)}).Delete(&entity.Friendships{}).Error
}
func (r *FriendshipsRepository) GetFriendshipByID(ctx context.Context, ID string) (*entity.Friendships, error) {
	var friendship *entity.Friendships
	err := r.db.WithContext(ctx).Model(&entity.Friendships{}).Where(&entity.Friendships{ID: uuid.MustParse(ID)}).First(&friendship).Error
	if err != nil {
		return nil, err
	}
	return friendship, nil
}

func (r *FriendshipsRepository) PanigationStatusFriendship(userID uuid.UUID, cursor string, limit int, status enum.StatusFriendship) (*dto.PaginationRes, error) {
	var data = []*entity.Friendships{}
	items := []string{
		"friendship_" + string(status) + "_cache_user_" + userID.String(),
		"friendship_" + string(status) + "_cache_nextcursor_user_" + userID.String(),
		"friendship_" + string(status) + "_cache_hasnext_user_" + userID.String(),
		"friendship_" + string(status) + "_cache_limit_user_" + userID.String(),
	}
	cachedData, nextCursor, hasNext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), items)
	if err != nil {
		return nil, err
	}
	if cachedData != nil {
		return &dto.PaginationRes{
			NextCursor: nextCursor,
			HasNext:    hasNext,
			Data:       cachedData,
			Limit:      limit,
		}, nil

	}
	querylimit := limit + 1
	query := r.db.Where(&entity.Friendships{Status: enum.StatusFriendship_Accepted}).
		Where("requester_id = ? OR recipient_id = ?", userID, userID).
		Order("created_at DESC ,id DESC").Limit(querylimit)
	if cursor != "" {
		time, datacursor, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where("(created_at > ?) OR (created_at = ? AND id > ?)", time, time, datacursor.ID)
		err = query.Find(data).Error
		if err != nil {
			return nil, err
		}
		hasnext := false
		if len(data) == querylimit {
			hasnext = true
			data = data[:limit]
		}
		lastdata := data[len(data)-1]
		nextcursor := utils.EncodeCursor(lastdata.Created_At, lastdata.ID)
		return &dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}, nil
	}
	err = query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	hasnext := false
	if len(data) == querylimit {
		hasnext = true
		data = data[:limit]
	}
	var nextcursor string
	if len(data) > 0 {
		lastdata := data[len(data)-1]
		nextcursor = utils.EncodeCursor(lastdata.Created_At, lastdata.ID)
	}
	var itemCache = map[string]any{
		"friendship_" + string(status) + "_cache_user_" + userID.String():            data,
		"friendship_" + string(status) + "_cache_nextcursor_user_" + userID.String(): nextcursor,
		"friendship_" + string(status) + "_cache_hasnext_user_" + userID.String():    hasnext,
		"friendship_" + string(status) + "_cache_limit_user_" + userID.String():      limit,
	}
	err = r.redisRepo.CustomizeSetCache(context.Background(), itemCache)
	if err != nil {
		return nil, err
	}
	return &dto.PaginationRes{
		NextCursor: nextcursor,
		HasNext:    hasnext,
		Data:       data,
		Limit:      limit,
	}, nil
}
