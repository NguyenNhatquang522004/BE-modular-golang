package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type DeleteAdCampainUsecase struct {
	accountRepo  IRepositoryPostgres.IAdAccountsRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	campaignRepo IRepositoryPostgres.IAdCampaignsRepository
}

func NewDeleteAdCampainUsecase(accountRepo IRepositoryPostgres.IAdAccountsRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, campaignRepo IRepositoryPostgres.IAdCampaignsRepository) *DeleteAdCampainUsecase {
	return &DeleteAdCampainUsecase{
		accountRepo:  accountRepo,
		pageRoleRepo: pageRoleRepo,
		campaignRepo: campaignRepo,
	}
}

func (u *DeleteAdCampainUsecase) Execute(ctx context.Context, req *req.DeleteAdCampainRequest) (*res.FailedAdCampainResponse, error) {
	// Implement the logic to delete an ad campaign here
	// This is a placeholder implementation and should be replaced with actual logic
	datacampain, err := u.campaignRepo.GetAdCampaignByIDDetail(ctx, req.CampaignID)
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID,
			UserActionID: req.UserActionID,
			AccountID:    req.UserActionID,
			ErrorMessage: "Failed to get ad campaign by ID: " + err.Error(),
		}, err
	}
	if datacampain == nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID,
			UserActionID: req.UserActionID,
			AccountID:    req.UserActionID,
			ErrorMessage: "Ad campaign not found",
		}, nil
	}
	if datacampain.AccountID.String() != req.UserActionID {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID,
			UserActionID: req.UserActionID,
			AccountID:    req.UserActionID,
			ErrorMessage: "User does not have permission to delete this ad campaign",
		}, nil
	}
	// Here you would call the repository method to delete the campaign, e.g.:
	err = u.campaignRepo.DeleteAdCampaign(ctx, req.CampaignID)
	if err != nil {
		return &res.FailedAdCampainResponse{
			CampaignID:   req.CampaignID,
			UserActionID: req.UserActionID,
			AccountID:    req.UserActionID,
			ErrorMessage: "Failed to delete ad campaign: " + err.Error(),
		}, err
	}
	return &res.FailedAdCampainResponse{
		CampaignID:   req.CampaignID,
		UserActionID: req.UserActionID,
		AccountID:    req.UserActionID,
		ErrorMessage: "Ad campaign deleted successfully",
	}, nil
}
