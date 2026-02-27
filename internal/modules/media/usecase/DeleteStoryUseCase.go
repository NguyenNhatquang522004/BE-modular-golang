package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type DeleteStoryUseCase struct {
	storyRepo IRepositoryMongodb.IStoryRepository
}

func NewDeleteStoryUseCase(storyRepo IRepositoryMongodb.IStoryRepository) *DeleteStoryUseCase {
	return &DeleteStoryUseCase{
		storyRepo: storyRepo,
	}
}

func (u *DeleteStoryUseCase) Execute(ctx context.Context, req *req.DeleteStoryRequest) (*res.FailedStoryResponse, error) {
	datastory, err := u.storyRepo.GetStoryByID(ctx, req.StoryId)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryId,
			UserID:       datastory.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	err = u.storyRepo.DeleteStory(ctx, req.StoryId)
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
