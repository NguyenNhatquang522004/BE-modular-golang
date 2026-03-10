package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type StoryStatsUseCase struct {
	events events.EventBus
}

func NewStoryStatsUseCase(events events.EventBus) *StoryStatsUseCase {
	return &StoryStatsUseCase{
		events: events,
	}
}

func (s *StoryStatsUseCase) Execute(ctx context.Context, req *req.StoryStatsRequest) (*res.FailedStoryResponse, error) {
	// Implement the logic to update story stats here
	// This might involve updating the database with the new stats for the story
	// You can also add any necessary validation or error handling as needed
	payload := &mediaEvent.StoryStatsPayload{
		StoryID:         req.StoryID,
		UserID:          req.UserID,
		Views:           req.Views,
		Like:            req.Like,
		Love:            req.Love,
		Wow:             req.Wow,
		Sad:             req.Sad,
		Angry:           req.Angry,
		ReplyCount:      req.ReplyCount,
		ViewsCount:      req.ViewsCount,
		InteractionType: req.InteractionType,
		PollOptionIndex: req.PollOptionIndex,
		Content:         req.Content,
		EventType:       req.EventType, // "increment" hoặc "decrement"
	}
	err := s.events.Publish(ctx, constants.TopicStoryStats.String(), req.UserID, req.EventType.String(), payload)
	if err != nil {
		return &res.FailedStoryResponse{
			StoryID:      req.StoryID,
			UserID:       req.UserID,
			ErrorMessage: err.Error(),
		}, err
	}
	// For now, we'll return a placeholder response
	return &res.FailedStoryResponse{
		StoryID:      req.StoryID,
		UserID:       req.UserID,
		ErrorMessage: "",
	}, nil
}
