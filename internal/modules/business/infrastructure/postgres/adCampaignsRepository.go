package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/postgresErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdCampaignsRepository struct {
	// Define the fields for the AdCampaignsRepository struct here
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

// Implement the methods for the AdCampaignsRepository struct here
func NewAdCampaignsRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *AdCampaignsRepository {
	return &AdCampaignsRepository{
		db:        db,
		redisRepo: redisRepo,
	}
}
func (r *AdCampaignsRepository) CreateAdCampaigns(ctx context.Context, campaign *entity.AdCampaign) error {
	return r.db.WithContext(ctx).Create(campaign).Error
}
func (r *AdCampaignsRepository) CreateBulkAdCampaigns(ctx context.Context, campaigns []*entity.AdCampaign) ([]*entity.AdCampaign, []*postgresErrors.AdCampaignsBulkError, error) {
	if len(campaigns) == 0 {
		return nil, nil, nil
	}
	var createdCampaigns []*entity.AdCampaign

	var bulkErrors []*postgresErrors.AdCampaignsBulkError
	for index, item := range campaigns {
		if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
			createdCampaigns = append(createdCampaigns, item)
			bulkErrors = append(bulkErrors, &postgresErrors.AdCampaignsBulkError{
				ID:        campaigns[index].ID.String(),
				AccountID: campaigns[index].AccountID.String(),
				Reason:    err.Error(),
			})
		}
	}
	return createdCampaigns, bulkErrors, nil
}
func (r *AdCampaignsRepository) GetAdCampaignByIDDetail(ctx context.Context, campaignID string) (*entity.AdCampaign, error) {
	var campaign *entity.AdCampaign
	if err := r.db.WithContext(ctx).Where(&entity.AdCampaign{ID: uuid.MustParse(campaignID)}).First(&campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}
func (r *AdCampaignsRepository) GetAdCampaign(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"campaign_cache_bulk_adcampaigns",
			"campaign_cache_bulk_nextcursor_adcampaigns",
			"campaign_cache_bulk_hasnext_adcampaigns",
			"campaign_cache_bulk_limit_adcampaigns",
		}
		data, nextcursor, hasnext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				Limit:      limit,
				NextCursor: nextcursor,
				HasNext:    hasnext,
				Data:       data,
			}, nil
		}
		return &dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}, nil
	}

	var campaigns []*entity.AdCampaign
	querylimit := limit + 1
	query := r.db.WithContext(ctx).Order("created_at DESC , id DESC").Limit(querylimit)
	if cursor != "" {
		createdAt, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", createdAt, createdAt, id)
	}
	err := query.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	var nextCursor string
	if len(campaigns) > 0 {
		lastCampaign := campaigns[len(campaigns)-1]
		nextCursor = utils.EncodeCursor(lastCampaign.CreatedAt, lastCampaign.ID)
	}
	hasNext := len(campaigns) > limit
	if hasNext {
		campaigns = campaigns[:limit]
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"campaign_cache_bulk_adcampaigns":            campaigns,
			"campaign_cache_bulk_nextcursor_adcampaigns": nextCursor,
			"campaign_cache_bulk_hasnext_adcampaigns":    hasNext,
			"campaign_cache_bulk_limit_adcampaigns":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       campaigns,
	}, nil
}
func (r *AdCampaignsRepository) GetAdCampaignsByAccountID(ctx context.Context, accountID string, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"campaign_cache_bulk_adcampaigns_accountID_" + accountID,
			"campaign_cache_bulk_nextcursor_adcampaigns_accountID_" + accountID,
			"campaign_cache_bulk_hasnext_adcampaigns_accountID_" + accountID,
			"campaign_cache_bulk_limit_adcampaigns_accountID_" + accountID,
		}
		data, nextcursor, hasnext, limit, err := r.redisRepo.CustomizeGetCache(context.Background(), itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				Limit:      limit,
				NextCursor: nextcursor,
				HasNext:    hasnext,
				Data:       data,
			}, nil
		}
		return &dto.PaginationRes{
			NextCursor: nextcursor,
			HasNext:    hasnext,
			Data:       data,
			Limit:      limit,
		}, nil
	}

	var campaigns []*entity.AdCampaign
	querylimit := limit + 1
	query := r.db.WithContext(ctx).Where(&entity.AdCampaign{AccountID: uuid.MustParse(accountID)}).Order("created_at DESC , id DESC").Limit(querylimit)
	if cursor != "" {
		createdAt, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", createdAt, createdAt, id)
	}
	err := query.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	var nextCursor string
	if len(campaigns) > 0 {
		lastCampaign := campaigns[len(campaigns)-1]
		nextCursor = utils.EncodeCursor(lastCampaign.CreatedAt, lastCampaign.ID)
	}
	hasNext := len(campaigns) > limit
	if hasNext {
		campaigns = campaigns[:limit]
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"campaign_cache_bulk_adcampaigns_accountID_" + accountID:            campaigns,
			"campaign_cache_bulk_nextcursor_adcampaigns_accountID_" + accountID: nextCursor,
			"campaign_cache_bulk_hasnext_adcampaigns_accountID_" + accountID:    hasNext,
			"campaign_cache_bulk_limit_adcampaigns_accountID_" + accountID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       campaigns,
	}, nil
}
func (r *AdCampaignsRepository) UpdateAdCampaign(ctx context.Context, campaign *entity.AdCampaign) error {
	return r.db.WithContext(ctx).Model(&entity.AdCampaign{}).Where(&entity.AdCampaign{ID: campaign.ID}).Updates(campaign).Error
}
func (r *AdCampaignsRepository) UpdateBulkAdCampaigns(ctx context.Context, campaigns []*entity.AdCampaign) ([]*entity.AdCampaign, []*postgresErrors.AdCampaignsBulkError, error) {
	if len(campaigns) == 0 {
		return nil, nil, nil
	}
	var updatedCampaigns []*entity.AdCampaign
	var bulkErrors []*postgresErrors.AdCampaignsBulkError
	for index, item := range campaigns {
		if err := r.db.WithContext(ctx).Model(&entity.AdCampaign{}).Where(&entity.AdCampaign{ID: item.ID}).Updates(item).Error; err != nil {
			updatedCampaigns = append(updatedCampaigns, item)
			bulkErrors = append(bulkErrors, &postgresErrors.AdCampaignsBulkError{
				ID:        campaigns[index].ID.String(),
				AccountID: campaigns[index].AccountID.String(),
				Reason:    err.Error(),
			})
		}
	}
	return updatedCampaigns, bulkErrors, nil
}
func (r *AdCampaignsRepository) DeleteAdCampaign(ctx context.Context, campaignID string) error {
	return r.db.WithContext(ctx).Where(&entity.AdCampaign{ID: uuid.MustParse(campaignID)}).Delete(&entity.AdCampaign{}).Error
}
func (r *AdCampaignsRepository) DeleteAdCampaignsBulk(ctx context.Context, campaignIDs []string) ([]*string, []*postgresErrors.AdCampaignsBulkError, error) {
	if len(campaignIDs) == 0 {
		return nil, nil, nil
	}
	var deletedCampaigns []*string
	var bulkErrors []*postgresErrors.AdCampaignsBulkError
	for _, campaignID := range campaignIDs {
		if err := r.db.WithContext(ctx).Where(&entity.AdCampaign{ID: uuid.MustParse(campaignID)}).Delete(&entity.AdCampaign{}).Error; err != nil {
			deletedCampaigns = append(deletedCampaigns, &campaignID)
		} else {
			bulkErrors = append(bulkErrors, &postgresErrors.AdCampaignsBulkError{
				ID:        campaignID,
				AccountID: "", // Không có AccountID trong trường hợp này, có thể để trống hoặc lấy từ DB nếu cần
				Reason:    err.Error(),
			})
		}
	}
	return deletedCampaigns, bulkErrors, nil
}
