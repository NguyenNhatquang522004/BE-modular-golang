package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type CreateAdCampainUsecase struct {
	accountRepo  IRepositoryPostgres.IAdAccountsRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	campaignRepo IRepositoryPostgres.IAdCampaignsRepository
}

func NewCreateAdCampainUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, campaignRepo IRepositoryPostgres.IAdCampaignsRepository) *CreateAdCampainUsecase {
	return &CreateAdCampainUsecase{
		accountRepo:  accountRepo,
		pageRoleRepo: pageRoleRepo,
		campaignRepo: campaignRepo,
	}
}

func (u *CreateAdCampainUsecase) Execute(ctx context.Context, req *req.CreateAdCampainRequest) (*res.FailedAdCampainResponse, error) {
	// Implement the logic to create an ad campaign here
	// This is a placeholder implementation and should be replaced with actual logic
	dataaccount, err := u.accountRepo.GetAdAccountsByOwnerUserIDDetail(ctx, req.AccountID.String())
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   "placeholder_campaign_id",
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "Failed to get ad accounts for user: " + err.Error(),
		}, err
	}
	if dataaccount == nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   "placeholder_campaign_id",
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "User does not have an ad account, cannot create campaign",
		}, nil
	}
	if dataaccount.OwnerUserID.String() != req.UserActionID {
		return &res.FailedAdCampainResponse{
			CampaignID:   "placeholder_campaign_id",
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "User does not have permission to create campaign for this ad account",
		}, nil
	}
	entityCampain := mapper.ToEntityAdCampaign(*req.AdCampaignReq)
	err = u.campaignRepo.CreateAdCampaigns(ctx, &entityCampain)
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   "placeholder_campaign_id",
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "Failed to create ad campaign: " + err.Error(),
		}, err
	}
	return &res.FailedAdCampainResponse{
		CampaignID:   "placeholder_campaign_id",
		UserActionID: req.UserActionID,
		AccountID:    req.AccountID.String(),
		ErrorMessage: "Failed to create ad campaign: not implemented",
	}, nil
}
