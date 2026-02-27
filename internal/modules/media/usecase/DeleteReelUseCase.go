package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type DeleteReelUseCase struct {
	reelRepo IRepositoryMongodb.IReelRepository
}

func NewDeleteReelUseCase(reelRepo IRepositoryMongodb.IReelRepository) *DeleteReelUseCase {
	return &DeleteReelUseCase{
		reelRepo: reelRepo,
	}
}
func (uc *DeleteReelUseCase) Execute(ctx context.Context, req *req.DeleteReelRequest) (*res.FailedReelResponse, error) {
	// Implement the logic for deleting a reel here
	data, err := uc.reelRepo.GetReelByID(ctx, req.ReelID)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       data.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	err = uc.reelRepo.DeleteReel(ctx, req.ReelID)
	if err != nil {
		return &res.FailedReelResponse{
			ReelID:       req.ReelID,
			UserID:       data.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedReelResponse{
		ReelID:       req.ReelID,
		UserID:       data.UserID,
		ErrorMessage: "",
	}, nil
	return &res.FailedReelResponse{
		ReelID:       req.ReelID,
		UserID:       data.UserID,
		ErrorMessage: "",
	}, nil
}
