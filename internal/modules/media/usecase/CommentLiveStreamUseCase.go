package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
)

type CommentLiveStreamUseCase struct {
	liveComment IRepositoryCassandra.ILiveCommentsRepository
	eventbus    events.EventBus
}

func NewCommentLiveStreamUseCase(liveComment IRepositoryCassandra.ILiveCommentsRepository, eventbus events.EventBus) *CommentLiveStreamUseCase {
	return &CommentLiveStreamUseCase{
		liveComment: liveComment,
		eventbus:    eventbus,
	}
}

func (c *CommentLiveStreamUseCase) Execute(ctx context.Context, reqa *req.CommentLiveStreamRequest) (*res.FailedLiveStreamResponse, error) {
	payload := &mediaEvent.LiveCommentPayload{
		StreamID:   reqa.StreamID,
		CreatedAt:  reqa.CreatedAt,
		CommentID:  reqa.CommentID,
		UserID:     reqa.UserID,
		UserBadges: reqa.UserBadges,
		Content:    reqa.Content,
		IsPinned:   reqa.IsPinned,
		EventType:  reqa.EventType,
	}
	err := c.eventbus.Publish(ctx, constants.TopicCommentLive.String(), reqa.StreamID, reqa.EventType.String(), payload)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: reqa.StreamID,
			UserID:        reqa.UserID,
			ErrorMessage:  err.Error(),
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: reqa.StreamID,
		UserID:        reqa.UserID,
		ErrorMessage:  "",
	}, nil
}
