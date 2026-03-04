package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type CreateadsUsecase struct {
	adsRepo IRepositoryPostgres.IAdsRepository
}

func NewCreateadsUsecase(adsRepo IRepositoryPostgres.IAdsRepository) *CreateadsUsecase {
	return &CreateadsUsecase{
		adsRepo: adsRepo,
	}
}
func (u *CreateadsUsecase) Execute(ctx context.Context, req *req.CreateAdsRequest) (*res.FailedAdResponse, error) {
	entity := mapper.ToEntityAds(*req.AdReq)
	err := u.adsRepo.CreateAd(ctx, &entity)
	if err != nil {
		return &res.FailedAdResponse{
			AdID:         "placeholder_ad_id",
			UserActionID: req.UserActionID,
			CampaignID:   req.AdReq.CampaignID.String(),
			ErrorMessage: "Failed to create ad: " + err.Error(),
		}, err
	}
	return &res.FailedAdResponse{
		AdID:         entity.ID.String(),
		UserActionID: req.UserActionID,
		CampaignID:   req.AdReq.CampaignID.String(),
		ErrorMessage: "",
	}, nil
}
