package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
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
	err := c.eventbus.Publish(ctx, constants.TopicCommentLive.String(), reqa.StreamID.String(), reqa.EventType.String(), reqa)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: reqa.StreamID.String(),
			UserID:        reqa.UserID.String(),
			ErrorMessage:  "Failed to publish comment event to event bus",
		}, err
	}
	counterReq := &req.CounterLiveStreamRequest{
		LiveSessionID: reqa.StreamID.String(),
		Comments:      1,
		Views:         0,
		EventType:     constants.Created, // Hoặc constants.Created tùy vào logic của bạn
	}
	err = c.eventbus.Publish(ctx, constants.TopicCounterLive.String(), reqa.StreamID.String(), reqa.EventType.String(), counterReq)
	if err != nil {
		return &res.FailedLiveStreamResponse{
			LiveSessionID: reqa.StreamID.String(),
			UserID:        reqa.UserID.String(),
			ErrorMessage:  "Failed to publish counter event to event bus",
		}, err
	}
	return &res.FailedLiveStreamResponse{
		LiveSessionID: reqa.StreamID.String(),
		UserID:        reqa.UserID.String(),
		ErrorMessage:  "Failed to add comment to live stream",
	}, nil
}
