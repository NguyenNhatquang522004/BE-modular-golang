package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepostitoryMongodb"
)

type CreateStoryUseCase struct {
	storyRepo IRepostitoryMongodb.IStoryRepository
}

func NewCreateStoryUseCase(storyRepo IRepostitoryMongodb.IStoryRepository) *CreateStoryUseCase {
	return &CreateStoryUseCase{
		storyRepo: storyRepo,
	}
}

func (uc *CreateStoryUseCase) Execute(ctx context.Context, req *req.CreateStoryRequest) (*res.FailedStoryResponse, error) {
	// Thực hiện logic tạo story ở đây
	entity, err := mapper.ToEntityStory(req.StoryReq)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      "",
			UserID:       req.StoryReq.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	err = uc.storyRepo.CreateStory(ctx, entity)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      "",
			UserID:       req.StoryReq.UserID,
			ErrorMessage: err.Error(),
		}, err
	}

	return &res.FailedStoryResponse{
		StoryID:      entity.ID.Hex(),
		UserID:       req.StoryReq.UserID,
		ErrorMessage: "",
	}, nil
}
