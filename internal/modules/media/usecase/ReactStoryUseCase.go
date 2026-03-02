package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ReactStoryUseCase struct {
	eventsbus events.EventBus
}

func NewReactStoryUseCase(eventsbus events.EventBus) *ReactStoryUseCase {
	return &ReactStoryUseCase{
		eventsbus: eventsbus,
	}
}
func (u *ReactStoryUseCase) Execute(ctx context.Context, req *req.ReactStoryRequest) (*res.FailedStoryResponse, error) {
	// Implement the logic to handle reacting to a story
	// This could involve validating the request, updating the story's reaction count, and returning an appropriate response
	payload := &mediaEvent.ReactStoryPayload{
		StoryId:         req.StoryId,
		UserId:          req.UserId,
		Avatar:          req.Avatar,
		Name:            req.Name,
		ViewedAt:        req.ViewedAt,
		InteractionType: req.InteractionType,
		ReactionCode:    req.ReactionCode,
		PollOptionIndex: req.PollOptionIndex,
		EventType:       req.EventType,
		Content:         req.Content,
	}
	err := u.eventsbus.Publish(ctx, constants.TopicReactStory.String(), req.StoryId, constants.Created.String(), payload)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryId,
			UserID:       req.UserId,
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedStoryResponse{
		StoryID: req.StoryId,
		UserID:  req.UserId,
	}, nil
}
