package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type DeleteadsUsecase struct {
	adsRepo IRepositoryPostgres.IAdsRepository
}

func NewDeleteadsUsecase(adsRepo IRepositoryPostgres.IAdsRepository) *DeleteadsUsecase {
	return &DeleteadsUsecase{
		adsRepo: adsRepo,
	}
}
func (u *DeleteadsUsecase) Execute(ctx context.Context, req *req.DeleteAdsRequest) (*res.FailedAdResponse, error) {
	err := u.adsRepo.DeleteAd(ctx, req.AdID)
	if err != nil {
		return &res.FailedAdResponse{
			AdID:         req.AdID,
			UserActionID: req.UserActionID,
			CampaignID:   " placeholder_campaign_id",
			ErrorMessage: "Failed to delete ad: " + err.Error(),
		}, err
	}
	return &res.FailedAdResponse{
		AdID:         req.AdID,
		UserActionID: req.UserActionID,
		CampaignID:   " placeholder_campaign_id",
		ErrorMessage: "",
	}, nil
}
