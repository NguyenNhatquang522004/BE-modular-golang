package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/postgresErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IAdsRepository interface {
	// Define the methods for the AdsRepository interface here
	CreateAd(ctx context.Context, ad *entity.Ad) error
	CreateAdsBulk(ctx context.Context, ads []*entity.Ad) ([]*entity.Ad, []*postgresErrors.AdsBulkError, error)
	GetAdByID(ctx context.Context, adID string) (*entity.Ad, error)
	GetBulkAdByID(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	GetAdsByCampaignID(ctx context.Context, campaignID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateAd(ctx context.Context, ad *entity.Ad) error
	UpdateAdsBulk(ctx context.Context, ads []*entity.Ad) ([]*entity.Ad, []*postgresErrors.AdsBulkError, error)
	DeleteAd(ctx context.Context, adID string) error
	DeleteAdsBulk(ctx context.Context, adIDs []string) ([]string, []*postgresErrors.AdsBulkError, error)
	DeleteAdsByCampaignID(ctx context.Context, campaignID string) error
	DeleteAdsByCampaignIDs(ctx context.Context, campaignIDs []string) ([]string, []*postgresErrors.AdsBulkError, error)
	DeleteAdsByTargetPostIDAndCampaignID(ctx context.Context, targetPostID string, campaignID string) error
}
