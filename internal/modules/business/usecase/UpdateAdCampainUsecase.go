package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type UpdateAdCampainUsecase struct {
	accountRepo  IRepositoryPostgres.IAdAccountsRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	campaignRepo IRepositoryPostgres.IAdCampaignsRepository
}

func NewUpdateAdCampainUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, campaignRepo IRepositoryPostgres.IAdCampaignsRepository) *UpdateAdCampainUsecase {
	return &UpdateAdCampainUsecase{
		accountRepo:  accountRepo,
		pageRoleRepo: pageRoleRepo,
		campaignRepo: campaignRepo,
	}
}

func (u *UpdateAdCampainUsecase) Execute(ctx context.Context, req *req.UpdateAdCampainRequest) (*res.FailedAdCampainResponse, error) {
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
	dataCampain, err := u.campaignRepo.GetAdCampaignByIDDetail(ctx, req.CampaignID.String())
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID.String(),
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "Failed to get ad campaign by ID: " + err.Error(),
		}, err
	}
	if dataCampain == nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID.String(),
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "Ad campaign not found",
		}, nil
	}
	if dataCampain.AccountID.String() != req.AccountID.String() {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID.String(),
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "User does not have permission to update this campaign",
		}, nil
	}
	mapper.UpdateToEntityAdCampaign(*req.AdCampaignReq, dataCampain)
	err = u.campaignRepo.UpdateAdCampaign(ctx, dataCampain)
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID.String(),
			UserActionID: req.UserActionID,
			AccountID:    req.AccountID.String(),
			ErrorMessage: "Failed to update ad campaign: " + err.Error(),
		}, err
	}

	return &res.FailedAdCampainResponse{
		CampaignID:   req.CampaignID.String(),
		UserActionID: req.UserActionID,
		AccountID:    req.AccountID.String(),
		ErrorMessage: "Ad campaign updated successfully",
	}, nil
}
