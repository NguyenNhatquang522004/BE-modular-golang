package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepostitoryMongodb"
)

type UpdateStoryUseCase struct {
	storyRepo IRepostitoryMongodb.IStoryRepository
}

func NewUpdateStoryUseCase(storyRepo IRepostitoryMongodb.IStoryRepository) *UpdateStoryUseCase {
	return &UpdateStoryUseCase{
		storyRepo: storyRepo,
	}
}

func (uc *UpdateStoryUseCase) Execute(ctx context.Context, req *req.UpdateStoryRequest) (*res.FailedStoryResponse, error) {
	// Thực hiện logic cập nhật story ở đây
	datastory, err := uc.storyRepo.GetStoryByID(ctx, req.StoryId)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryId,
			UserID:       datastory.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	mapper.UpdateToEntityStory(req.UpdateStoryReq, datastory)
	err = uc.storyRepo.UpdateStory(ctx, datastory)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryId,
			UserID:       datastory.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedStoryResponse{
		StoryID:      req.StoryId,
		UserID:       datastory.UserID,
		ErrorMessage: "",
	}, nil
}
