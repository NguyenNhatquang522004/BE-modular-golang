package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type CommentCountUsecase struct {
	commentRepo IRepositoryMongoDB.ICommentRepository
	events      events.EventBus
}

func NewCommentCountUsecase(commentRepo IRepositoryMongoDB.ICommentRepository, events events.EventBus) *CommentCountUsecase {
	return &CommentCountUsecase{
		commentRepo: commentRepo,
		events:      events,
	}
}

func (c *CommentCountUsecase) Execute(ctx context.Context, req *req.CommentCountRequest) (*res.FailedCommentCountResponse, error) {
	err := c.events.Publish(ctx, constants.TopicCounterComment.String(), req.CommentID, req.EventType.String(), req)
	if err != nil {
		return &res.FailedCommentCountResponse{
			CommentID:    req.CommentID,
			ReplyCount:   req.ReplyCount,
			MentionCount: req.MentionCount,
			ReportCount:  req.ReportCount,
			EventType:    req.EventType.String(),
			ErrorMessage: err.Error(),
		}, err
	}
	return &res.FailedCommentCountResponse{
		CommentID:    req.CommentID,
		ReplyCount:   req.ReplyCount,
		MentionCount: req.MentionCount,
		ReportCount:  req.ReportCount,
		EventType:    req.EventType.String(),
		ErrorMessage: "",
	}, nil
}
