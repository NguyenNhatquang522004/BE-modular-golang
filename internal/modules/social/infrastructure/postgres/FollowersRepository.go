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

func (r *FollowersRepository) CreateFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID) error {
	return r.db.Create(&entity.Followers{
		Follower_UserID: followerUserID,
		Followed_UserID: followedUserID,
		IsMuted:         false,
	}).Error
}

func (r *FollowersRepository) DeleteSoftFollowUser(follower uuid.UUID) error {
	data, err := r.GetFollowerByID(follower)
	if err != nil {
		return err
	}
	err = r.db.Model(&data).Update("deleted_at", gorm.DeletedAt{Time: data.UpdatedAt, Valid: true}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *FollowersRepository) DeleteBatchSoftFollowUser(followerID uuid.UUID, followedUserID uuid.UUID) error {
	// Implementation here
	err := r.db.Where(&entity.Followers{Follower_UserID: followerID, Followed_UserID: followedUserID}).
		Or(&entity.Followers{Follower_UserID: followedUserID, Followed_UserID: followerID}).
		Update("deleted_at", gorm.Expr("updated_at")).
		Error

	if err != nil {
		return err
	}
	return nil
}
func (r *FollowersRepository) DeleteBatchHardFollowUser(followerID uuid.UUID, followedUserID uuid.UUID) error {
	// Implementation here
	err := r.db.Where(&entity.Followers{Follower_UserID: followerID, Followed_UserID: followedUserID}).
		Or(&entity.Followers{Follower_UserID: followedUserID, Followed_UserID: followerID}).
		Delete(&entity.Followers{}).
		Error

	if err != nil {
		return err
	}
	return nil
}
func (r *FollowersRepository) DeleteHardFollowUser(follower uuid.UUID) error {
	data, err := r.GetFollowerByID(follower)
	if err != nil {
		return err
	}
	err = r.db.Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FollowersRepository) UpdatateMuteFollowUser(followerUserID uuid.UUID, followedUserID uuid.UUID, isMuted bool) error {
	data, err := r.GetFollowerBybidirectional(followerUserID, followedUserID)
	if err != nil {
		return err
	}
	data.IsMuted = isMuted
	r.db.Save(data)
	return nil
}

func (r *FollowersRepository) PaginationFollowers(FollowerUserID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implementation here
	var followers = []*entity.Followers{}
	queryLimit := limit + 1
	items := []string{
		"follower_cache_user_" + FollowerUserID.String(),
		"follower_cache_nextcursor_user_" + FollowerUserID.String(),
		"follower_cache_hasnext_user_" + FollowerUserID.String(),
		"follower_cache_limit_user_" + FollowerUserID.String(),
	}
	data, cursor, hasNextcache, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), items)
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

	query := r.db.Where(&entity.Followers{Follower_UserID: FollowerUserID}).Order("created_at DESC ,id DESC").Limit(queryLimit)
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
		"follower_cache_user_" + FollowerUserID.String():            followers,
		"follower_cache_nextcursor_user_" + FollowerUserID.String(): nextCursor,
		"follower_cache_hasnext_user_" + FollowerUserID.String():    hasNext,
		"follower_cache_limit_user_" + FollowerUserID.String():      limit,
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
func (r *FollowersRepository) PaginationFolloweds(FollowedUserID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implementation here
	queryLimit := limit + 1
	var followeds = []*entity.Followers{}
	items := []string{
		"followed_cache_user_" + FollowedUserID.String(),
		"followed_cache_nextcursor_user_" + FollowedUserID.String(),
		"followed_cache_hasnext_user_" + FollowedUserID.String(),
		"followed_cache_limit_user_" + FollowedUserID.String(),
	}
	data, cursor, hasNextcache, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), items)
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
	query := r.db.Where(&entity.Followers{Followed_UserID: FollowedUserID}).Order("created_at DESC ,id DESC  ").Limit(queryLimit)
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
		"followed_cache_user_" + FollowedUserID.String():            followeds,
		"followed_cache_nextcursor_user_" + FollowedUserID.String(): nextCursor,
		"followed_cache_hasnext_user_" + FollowedUserID.String():    hasNext,
		"followed_cache_limit_user_" + FollowedUserID.String():      limit,
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
func (r *FollowersRepository) GetFollowerByID(follower uuid.UUID) (*entity.Followers, error) {
	// Implementation her
	data := &entity.Followers{}
	err := r.db.Where(&entity.Followers{ID: follower}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *FollowersRepository) GetFollowerBybidirectional(followerUserID uuid.UUID, followedUserID uuid.UUID) (*entity.Followers, error) {
	// Implementation here
	data := &entity.Followers{}
	err := r.db.Where(&entity.Followers{Follower_UserID: followerUserID, Followed_UserID: followedUserID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (r *FollowersRepository) GetFollowerIndiscriminate(followerUserID uuid.UUID, followedUserID uuid.UUID) (*[]entity.Followers, error) {
	data := &[]entity.Followers{}
	err := r.db.Where(&entity.Followers{Follower_UserID: followerUserID, Followed_UserID: followedUserID}).Or(&entity.Followers{Follower_UserID: followedUserID, Followed_UserID: followerUserID}).Find(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
