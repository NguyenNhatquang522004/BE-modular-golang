package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowersRepository struct {
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

// Implement the methods defined in IFollowersRepository interface
func NewFollowersRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *FollowersRepository {
	return &FollowersRepository{
		db:        db,
		redisRepo: redisRepo,
	}
}

func (r *FollowersRepository) CreateFollower(ctx context.Context, req *entity.Followers) error {
	return r.db.WithContext(ctx).Model(&entity.Followers{}).Create(req).Error
}
func (r *FollowersRepository) UpdateFollower(ctx context.Context, req *entity.Followers) error {
	return r.db.WithContext(ctx).Model(&entity.Followers{}).Where(&entity.Followers{ID: req.ID}).Updates(req).Error
}
func (r *FollowersRepository) DeleteFollower(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Model(&entity.Followers{}).Where(&entity.Followers{ID: uuid.MustParse(ID)}).Delete(&entity.Followers{}).Error
}
func (r *FollowersRepository) DeleteFollowerByUserID(ctx context.Context, FollowerUserID string, FollowedUserID string) error {
	return r.db.WithContext(ctx).Model(&entity.Followers{}).Where(&entity.Followers{Follower_UserID: uuid.MustParse(FollowerUserID), Followed_UserID: uuid.MustParse(FollowedUserID)}).Or(&entity.Followers{Follower_UserID: uuid.MustParse(FollowedUserID), Followed_UserID: uuid.MustParse(FollowerUserID)}).Delete(&entity.Followers{}).Error
}
func (r *FollowersRepository) GetFollowerByID(ctx context.Context, ID string) (*entity.Followers, error) {
	var follower *entity.Followers
	err := r.db.WithContext(ctx).Model(&entity.Followers{}).Where(&entity.Followers{ID: uuid.MustParse(ID)}).First(&follower).Error
	if err != nil {
		return nil, err
	}
	return follower, nil
}
func (r *FollowersRepository) GetFollowerByUserIDs(ctx context.Context, FollowerUserID string, FollowedUserID string) (*entity.Followers, error) {
	var follower *entity.Followers
	err := r.db.WithContext(ctx).Model(&entity.Followers{}).Where(&entity.Followers{Follower_UserID: uuid.MustParse(FollowerUserID), Followed_UserID: uuid.MustParse(FollowedUserID)}).First(&follower).Error
	if err != nil {
		return nil, err
	}
	return follower, nil
}
func (r *FollowersRepository) PaginationFollowers(ctx context.Context, FollowerUserID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implementation here
	var followers = []*entity.Followers{}
	queryLimit := limit + 1
	items := []string{
		"follower_cache_user_" + FollowerUserID,
		"follower_cache_nextcursor_user_" + FollowerUserID,
		"follower_cache_hasnext_user_" + FollowerUserID,
		"follower_cache_limit_user_" + FollowerUserID,
	}
	data, cursor, hasNextcache, limit, err := r.redisRepo.CustomizeGetCache(ctx, items)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return &dto.PaginationRes{
			NextCursor: cursor,
			HasNext:    hasNextcache,
			Data:       data,
			Limit:      limit,
		}, nil
	}

	query := r.db.Where(&entity.Followers{Follower_UserID: uuid.MustParse(FollowerUserID)}).Order("created_at DESC ,id DESC").Limit(queryLimit)
	if cursor != "" {
		time, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where(" (created_at < ?) OR (created_at = ? AND id > ?)", time, time, id)
		err = query.Find(followers).Error
		if err != nil {
			return nil, err
		}
		var hasNext = false
		if len(followers) > limit {
			hasNext = true
			followers = followers[:limit]
		}

		lastFollower := followers[len(followers)-1]
		nextCursor := utils.EncodeCursor(lastFollower.CreatedAt, lastFollower.Followed_UserID)
		return &dto.PaginationRes{
			NextCursor: nextCursor,
			HasNext:    hasNext,
			Data:       followers,
			Limit:      limit,
		}, nil

	}
	err = query.Find(followers).Error
	if err != nil {
		return nil, err
	}
	var hasNext = false
	if len(followers) > limit {
		hasNext = true
		followers = followers[:limit]
	}

	var nextCursor string
	if len(followers) > 0 {
		lastFollower := followers[len(followers)-1]
		nextCursor = utils.EncodeCursor(lastFollower.CreatedAt, lastFollower.Followed_UserID)
	}
	itemsSet := map[string]any{
		"follower_cache_user_" + FollowerUserID:            followers,
		"follower_cache_nextcursor_user_" + FollowerUserID: nextCursor,
		"follower_cache_hasnext_user_" + FollowerUserID:    hasNext,
		"follower_cache_limit_user_" + FollowerUserID:      limit,
	}
	err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
	if err != nil {
		return nil, err
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       followers,
		Limit:      limit,
	}, nil
}
func (r *FollowersRepository) PaginationFolloweds(ctx context.Context, FollowedUserID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implementation here
	queryLimit := limit + 1
	var followeds = []*entity.Followers{}
	items := []string{
		"followed_cache_user_" + FollowedUserID,
		"followed_cache_nextcursor_user_" + FollowedUserID,
		"followed_cache_hasnext_user_" + FollowedUserID,
		"followed_cache_limit_user_" + FollowedUserID,
	}
	data, cursor, hasNextcache, limit, err := r.redisRepo.CustomizeGetCache(ctx, items)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return &dto.PaginationRes{
			NextCursor: cursor,
			HasNext:    hasNextcache,
			Data:       data,
			Limit:      limit,
		}, nil
	}
	query := r.db.Where(&entity.Followers{Followed_UserID: uuid.MustParse(FollowedUserID)}).Order("created_at DESC ,id DESC  ").Limit(queryLimit)
	if cursor != "" {
		time, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where(" (created_at < ?) OR (created_at = ? AND id > ?)", time, time, id)
		err = query.Find(followeds).Error
		if err != nil {
			return nil, err
		}
		var hasNext = false
		if len(followeds) > limit {
			hasNext = true
			followeds = followeds[:limit]
		}

		lastFollower := followeds[len(followeds)-1]
		nextCursor := utils.EncodeCursor(lastFollower.CreatedAt, lastFollower.ID)
		return &dto.PaginationRes{
			NextCursor: nextCursor,
			HasNext:    hasNext,
			Data:       followeds,
			Limit:      limit,
		}, nil

	}

	err = query.Find(followeds).Error
	if err != nil {
		return nil, err
	}
	var hasNext = false
	if len(followeds) > limit {
		hasNext = true
		followeds = followeds[:limit]
	}

	var nextCursor string
	if len(followeds) > 0 {
		lastFollower := followeds[len(followeds)-1]
		nextCursor = utils.EncodeCursor(lastFollower.CreatedAt, lastFollower.ID)
	}
	itemsSet := map[string]any{
		"followed_cache_user_" + FollowedUserID:            followeds,
		"followed_cache_nextcursor_user_" + FollowedUserID: nextCursor,
		"followed_cache_hasnext_user_" + FollowedUserID:    hasNext,
		"followed_cache_limit_user_" + FollowedUserID:      limit,
	}
	err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
	if err != nil {
		return nil, err
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       followeds,
		Limit:      limit,
	}, nil
}
