package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type UpdateReelUseCase struct {
	reelRepo IRepositoryMongodb.IReelRepository
}

func NewUpdateReelUseCase(reelRepo IRepositoryMongodb.IReelRepository) *UpdateReelUseCase {
	return &UpdateReelUseCase{
		reelRepo: reelRepo,
	}
}

func (u *UpdateReelUseCase) Execute(ctx context.Context, req *req.UpdateReelRequest) (*res.FailedReelResponse, error) {
	// Thực hiện logic cập nhật reel ở đây
	dataReel, err := u.reelRepo.GetReelByID(ctx, req.ReelID)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       dataReel.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	// Cập nhật các trường của dataReel dựa trên req.UpdateReelReq
	mapper.UpdateToEntityReel(req.UpdateReelReq, dataReel)
	err = u.reelRepo.UpdateReel(ctx, dataReel)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       dataReel.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedReelResponse{
		ReelID:       req.ReelID,
		UserID:       dataReel.UserID,
		ErrorMessage: "",
	}, nil
}
