package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ViewCountStoryUseCase struct {
	eventBus events.EventBus
}

func NewViewCountStoryUseCase(eventBus events.EventBus) *ViewCountStoryUseCase {
	return &ViewCountStoryUseCase{
		eventBus: eventBus,
	}
}

func (uc *ViewCountStoryUseCase) Execute(ctx context.Context, req *req.ViewCountStoryRequest) (*res.FailedStoryResponse, error) {
	err := uc.eventBus.Publish(ctx, constants.TopicViewCountStory.String(), req.StoryId, constants.Created.String(), req)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryId,
			UserID:       req.UserId,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedStoryResponse{
		StoryID:      req.StoryId,
		UserID:       req.UserId,
		ErrorMessage: "",
	}, nil
}
