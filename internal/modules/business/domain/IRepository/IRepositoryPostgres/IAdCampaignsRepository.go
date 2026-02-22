package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/postgresErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
)

type IAdCampaignsRepository interface {
	CreateAdCampaigns(ctx context.Context, campaigns []*entity.AdCampaign) error
	CreateBulkAdCampaigns(ctx context.Context, campaigns []*entity.AdCampaign) ([]*entity.AdCampaign, []*postgresErrors.AdCampaignsBulkError, error)
	GetAdCampaignByIDDetail(ctx context.Context, campaignID string) (*entity.AdCampaign, error)
	GetAdCampaign(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	GetAdCampaignsByAccountID(ctx context.Context, accountID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateAdCampaign(ctx context.Context, campaign *entity.AdCampaign) error
	UpdateBulkAdCampaigns(ctx context.Context, campaigns []*entity.AdCampaign) ([]*entity.AdCampaign, []*postgresErrors.AdCampaignsBulkError, error)
	DeleteAdCampaign(ctx context.Context, campaignID string) error
	DeleteAdCampaignsBulk(ctx context.Context, campaignIDs []string) ([]*string, []*postgresErrors.AdCampaignsBulkError, error)
}
