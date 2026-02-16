package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
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

func (r *FriendshipsRepository) CreateFriendship(Requester_ID uuid.UUID, Recipient_ID uuid.UUID) error {

	return r.db.Create(&entity.Friendships{
		Requester_ID: Requester_ID,
		Recipient_ID: Recipient_ID,
		Status:       enum.StatusFriendship_Pending,
	}).Error
}

func (r *FriendshipsRepository) UpdateFriendshipStatus(friendshipID uuid.UUID, status enum.StatusFriendship) error {
	data, err := r.GetFriendshipByUserIDs(friendshipID)
	if err != nil {
		return err
	}
	data.Status = status
	r.db.Save(data)
	return nil
}

func (r *FriendshipsRepository) DeleteHardFriendship(friendshipID uuid.UUID) error {
	data, err := r.GetFriendshipByUserIDs(friendshipID)
	if err != nil {
		return err
	}
	err = r.db.Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FriendshipsRepository) DeleteSoftFriendship(friendshipID uuid.UUID) error {
	data, err := r.GetFriendshipByUserIDs(friendshipID)
	if err != nil {
		return err
	}
	err = r.db.Model(&data).Update("deleted_at", gorm.DeletedAt{Time: data.Updated_At, Valid: true}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *FriendshipsRepository) PanigationAcceptedFriendship(userID uuid.UUID, cursor string, limit int) (*response.Response, error) {
	var data = []*entity.Friendships{}
	items := []string{
		"friendship_AC_cache_user_" + userID.String(),
		"friendship_AC_cache_nextcursor_user_" + userID.String(),
		"friendship_AC_cache_hasnext_user_" + userID.String(),
		"friendship_AC_cache_limit_user_" + userID.String(),
	}
	cachedData, nextCursor, hasNext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), items)
	if err != nil {
		return nil, err
	}
	if cachedData != nil {
		return response.NewResponse(response.WithData(&dto.PaginationRes{
			NextCursor: nextCursor,
			HasNext:    hasNext,
			Data:       cachedData,
			Limit:      limit,
		}), response.WithMessage(""), response.WithStatus("")), nil

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
		return response.NewResponse(response.WithData(&dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}), response.WithMessage(""), response.WithStatus("")), nil
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
		"friendship_AC_cache_user_" + userID.String():            data,
		"friendship_AC_cache_nextcursor_user_" + userID.String(): nextcursor,
		"friendship_AC_cache_hasnext_user_" + userID.String():    hasnext,
		"friendship_AC_cache_limit_user_" + userID.String():      limit,
	}
	err = r.redisRepo.CustomizeSetCache(context.Background(), itemCache)
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(&dto.PaginationRes{
		NextCursor: nextcursor,
		HasNext:    hasnext,
		Data:       data,
		Limit:      limit,
	}), response.WithMessage(""), response.WithStatus("")), nil
}
func (r *FriendshipsRepository) PanigationPendingFriendship(userID uuid.UUID, cursor string, limit int) (*dto.PaginationRes, error) {
	var data = []*entity.Friendships{}
	ctx := context.Background()
	items := []string{
		"friendship_PD_cache_user_" + userID.String(),
		"friendship_PD_cache_nextcursor_user_" + userID.String(),
		"friendship_PD_cache_hasnext_user_" + userID.String(),
		"friendship_PD_cache_limit_user_" + userID.String(),
	}
	cachedData, nextCursor, hasNext, limit, err := r.redisRepo.CustomizeGetCache(ctx, items)
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
	query := r.db.Where(&entity.Friendships{Status: enum.StatusFriendship_Pending}).
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

	err = query.Find(data).Error
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
	itemCache := map[string]any{
		"friendship_PD_cache_user_" + userID.String():            data,
		"friendship_PD_cache_nextcursor_user_" + userID.String(): nextcursor,
		"friendship_PD_cache_hasnext_user_" + userID.String():    hasnext,
		"friendship_PD_cache_limit_user_" + userID.String():      limit,
	}
	err = r.redisRepo.CustomizeSetCache(ctx, itemCache)
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
func (r *FriendshipsRepository) GetFriendshipByUserIDs(friendshipID uuid.UUID) (*entity.Friendships, error) {
	var data = &entity.Friendships{}
	err := r.db.Where(&entity.Friendships{ID: friendshipID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (r *FriendshipsRepository) GetFriendshipBybidirectional(requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error) {
	data := &entity.Friendships{}
	err := r.db.Where(&entity.Friendships{Requester_ID: requesterID, Recipient_ID: recipientID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (r *FriendshipsRepository) GetFriendshipIndiscriminate(requesterID uuid.UUID, recipientID uuid.UUID) (*entity.Friendships, error) {
	data := &entity.Friendships{}
	err := r.db.Where(&entity.Friendships{Requester_ID: requesterID, Recipient_ID: recipientID}).Or(&entity.Friendships{Requester_ID: recipientID, Recipient_ID: requesterID}).First(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
