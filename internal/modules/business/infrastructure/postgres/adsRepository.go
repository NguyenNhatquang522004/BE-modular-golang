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

type AdsRepository struct {
	// Define the fields for the AdRepository struct here
	db        *gorm.DB
	redisRepo IRepositoryShare.IRedis
}

func NewAdsRepository(db *gorm.DB, redisRepo IRepositoryShare.IRedis) *AdsRepository {
	return &AdsRepository{
		db:        db,
		redisRepo: redisRepo,
	}
}
func (r *AdsRepository) CreateAd(ctx context.Context, ad *entity.Ad) error {
	return r.db.WithContext(ctx).Create(ad).Error
}
func (r *AdsRepository) CreateAdsBulk(ctx context.Context, ads []*entity.Ad) ([]*entity.Ad, []*postgresErrors.AdsBulkError, error) {
	if len(ads) == 0 {
		return nil, nil, nil
	}
	var createdAds []*entity.Ad

	var bulkErrors []*postgresErrors.AdsBulkError
	for index, item := range ads {
		if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
			createdAds = append(createdAds, item)
			bulkErrors = append(bulkErrors, &postgresErrors.AdsBulkError{
				ID:           ads[index].ID.String(),
				CampaignID:   ads[index].CampaignID.String(),
				TargetPostID: ads[index].TargetPostID,
				Reason:       err.Error(),
			})
		}
	}
	return createdAds, bulkErrors, nil
}
func (r *AdsRepository) GetAdByID(ctx context.Context, adID string) (*entity.Ad, error) {
	var ad *entity.Ad
	if err := r.db.WithContext(ctx).Where("id = ?", adID).First(&ad).Error; err != nil {
		return nil, err
	}
	return ad, nil
}
func (r *AdsRepository) GetBulkAdByID(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"ads_cache_bulk_ads",
			"ads_cache_bulk_nextcursor_ads",
			"ads_cache_bulk_hasnext_ads",
			"ads_cache_bulk_limit_ads",
		}
		cachedData, nextcursor, hasNext, limit, err := r.redisRepo.CustomizeGetCache(ctx, itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				NextCursor: nextcursor,
				HasNext:    hasNext,
				Data:       cachedData,
				Limit:      limit,
			}, nil
		}
	}
	var ads []*entity.Ad
	querylimit := limit + 1
	query := r.db.Order("created_at DESC, id DESC").Limit(querylimit)
	if cursor != "" {
		time, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where(" (created_at < ?) OR (created_at = ? AND id > ?)", time, time, id)
	}
	err := query.Find(&ads).Error
	if err != nil {
		return nil, err
	}
	var hasNext = false
	if len(ads) > limit {
		hasNext = true
		ads = ads[:limit]
	}
	var nextCursor string
	if len(ads) > 0 {
		lastAd := ads[len(ads)-1]
		nextCursor = utils.EncodeCursor(lastAd.CreatedAt, lastAd.ID)
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"ads_cache_bulk_ads":            ads,
			"ads_cache_bulk_nextcursor_ads": nextCursor,
			"ads_cache_bulk_hasnext_ads":    hasNext,
			"ads_cache_bulk_limit_ads":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       ads,
		Limit:      limit,
	}, nil
}
func (r *AdsRepository) GetAdsByCampaignID(ctx context.Context, campaignID string, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		itemsSet := []string{
			"ads_campaign_cache_bulk_ads_campaignID_" + campaignID,
			"ads_campaign_cache_bulk_nextcursor_ads_campaignID_" + campaignID,
			"ads_campaign_cache_bulk_hasnext_ads_campaignID_" + campaignID,
			"ads_campaign_cache_bulk_limit_ads_campaignID_" + campaignID,
		}
		cachedData, nextcursor, hasNext, limit, err := r.redisRepo.CustomizeGetCache(ctx, itemsSet)
		if err == nil {
			return &dto.PaginationRes{
				NextCursor: nextcursor,
				HasNext:    hasNext,
				Data:       cachedData,
				Limit:      limit,
			}, nil
		}
	}
	var ads []*entity.Ad
	querylimit := limit + 1
	query := r.db.Where(&entity.Ad{CampaignID: uuid.MustParse(campaignID)}).Order("created_at DESC, id DESC").Limit(querylimit)
	if cursor != "" {
		time, id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}
		query = query.Where(" (created_at < ?) OR (created_at = ? AND id > ?)", time, time, id)
	}
	err := query.Find(&ads).Error
	if err != nil {
		return nil, err
	}
	var hasNext = false
	if len(ads) > limit {
		hasNext = true
		ads = ads[:limit]
	}
	var nextCursor string
	if len(ads) > 0 {
		lastAd := ads[len(ads)-1]
		nextCursor = utils.EncodeCursor(lastAd.CreatedAt, lastAd.ID)
	}
	if cursor == "" {
		itemsSet := map[string]any{
			"ads_campaign_cache_bulk_ads_campaignID_" + campaignID:            ads,
			"ads_campaign_cache_bulk_nextcursor_ads_campaignID_" + campaignID: nextCursor,
			"ads_campaign_cache_bulk_hasnext_ads_campaignID_" + campaignID:    hasNext,
			"ads_campaign_cache_bulk_limit_ads_campaignID_" + campaignID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(context.Background(), itemsSet)
		if err != nil {
			return nil, err
		}
	}
	return &dto.PaginationRes{
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Data:       ads,
		Limit:      limit,
	}, nil
}
func (r *AdsRepository) UpdateAd(ctx context.Context, ad *entity.Ad) error {
	return r.db.WithContext(ctx).Save(ad).Error
}
func (r *AdsRepository) UpdateAdsBulk(ctx context.Context, ads []*entity.Ad) ([]*entity.Ad, []*postgresErrors.AdsBulkError, error) {
	if len(ads) == 0 {
		return nil, nil, nil
	}
	var updatedAds []*entity.Ad

	var bulkErrors []*postgresErrors.AdsBulkError
	for index, item := range ads {
		if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
			updatedAds = append(updatedAds, item)
			bulkErrors = append(bulkErrors, &postgresErrors.AdsBulkError{
				ID:           ads[index].ID.String(),
				CampaignID:   ads[index].CampaignID.String(),
				TargetPostID: ads[index].TargetPostID,
				Reason:       err.Error(),
			})
		}
	}
	return updatedAds, bulkErrors, nil
}
func (r *AdsRepository) DeleteAd(ctx context.Context, adID string) error {
	return r.db.WithContext(ctx).Where(&entity.Ad{ID: uuid.MustParse(adID)}).Delete(&entity.Ad{}).Error
}
func (r *AdsRepository) DeleteAdsBulk(ctx context.Context, adIDs []string) ([]string, []*postgresErrors.AdsBulkError, error) {
	if len(adIDs) == 0 {
		return nil, nil, nil
	}
	var deletedAds []string
	var bulkErrors []*postgresErrors.AdsBulkError
	for _, adID := range adIDs {
		if err := r.db.WithContext(ctx).Where(&entity.Ad{ID: uuid.MustParse(adID)}).Delete(&entity.Ad{}).Error; err != nil {
			deletedAds = append(deletedAds, adID)
		} else {
			bulkErrors = append(bulkErrors, &postgresErrors.AdsBulkError{
				ID:     adID,
				Reason: err.Error(),
			})
		}
	}
	return deletedAds, bulkErrors, nil
}
func (r *AdsRepository) DeleteAdsByCampaignID(ctx context.Context, campaignID string) error {
	return r.db.WithContext(ctx).Where(&entity.Ad{CampaignID: uuid.MustParse(campaignID)}).Delete(&entity.Ad{}).Error
}
func (r *AdsRepository) DeleteAdsByCampaignIDs(ctx context.Context, campaignIDs []string) ([]string, []*postgresErrors.AdsBulkError, error) {
	if len(campaignIDs) == 0 {
		return nil, nil, nil
	}
	var deletedCampaignIDs []string
	var bulkErrors []*postgresErrors.AdsBulkError
	for _, campaignID := range campaignIDs {
		if err := r.db.WithContext(ctx).Where(&entity.Ad{CampaignID: uuid.MustParse(campaignID)}).Delete(&entity.Ad{}).Error; err != nil {
			deletedCampaignIDs = append(deletedCampaignIDs, campaignID)
		} else {
			bulkErrors = append(bulkErrors, &postgresErrors.AdsBulkError{
				ID:     campaignID,
				Reason: err.Error(),
			})
		}
	}
	return deletedCampaignIDs, bulkErrors, nil
}
func (r *AdsRepository) DeleteAdsByTargetPostIDAndCampaignID(ctx context.Context, targetPostID string, campaignID string) error {
	return r.db.WithContext(ctx).Where(&entity.Ad{TargetPostID: targetPostID, CampaignID: uuid.MustParse(campaignID)}).Delete(&entity.Ad{}).Error
}
