package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
)

type UpdateadsUsecase struct {
	adsRepo IRepositoryPostgres.IAdsRepository
}

func NewUpdateadsUsecase(adsRepo IRepositoryPostgres.IAdsRepository) *UpdateadsUsecase {
	return &UpdateadsUsecase{
		adsRepo: adsRepo,
	}
}
func (u *UpdateadsUsecase) Execute(ctx context.Context, req *req.UpdateAdsRequest) (*res.FailedAdResponse, error) {
	// Lấy ad hiện tại từ DB để cập nhật
	data, err := u.adsRepo.GetAdByID(ctx, req.AdReq.AdID.String())
	if err != nil {
		return &res.FailedAdResponse{
			AdID:         req.AdReq.AdID.String(),
			UserActionID: req.UserActionID,
			CampaignID:   req.AdReq.CampaignID.String(),
			ErrorMessage: "Failed to retrieve ad for update: " + err.Error(),
		}, err
	}
	if data == nil {
		return &res.FailedAdResponse{
			AdID:         req.AdReq.AdID.String(),
			UserActionID: req.UserActionID,
			CampaignID:   req.AdReq.CampaignID.String(),
			ErrorMessage: "Ad not found for update",
		}, nil
	}

	// Cập nhật các trường có giá trị hợp lệ từ request vào entity hiện tại
	mapper.UpdateToEntityAds(*req.AdReq, data)

	err = u.adsRepo.UpdateAd(ctx, data)
	if err != nil {
		return &res.FailedAdResponse{
			AdID:         req.AdReq.AdID.String(),
			UserActionID: req.UserActionID,
			CampaignID:   req.AdReq.CampaignID.String(),
			ErrorMessage: "Failed to update ad: " + err.Error(),
		}, err
	}
	return &res.FailedAdResponse{
		AdID:         data.ID.String(),
		UserActionID: req.UserActionID,
		CampaignID:   req.AdReq.CampaignID.String(),
		ErrorMessage: "",
	}, nil
}
