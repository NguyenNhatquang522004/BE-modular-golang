package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type DeleteCommentUsecase struct {
	eventBus           events.EventBus
	commentRepo        IRepositoryMongoDB.ICommentRepository
	commentEditLogRepo IRepositoryMongoDB.ICommentEditLogsRepository
}

func NewDeleteCommentUsecase(eventBus events.EventBus, commentRepo IRepositoryMongoDB.ICommentRepository, commentEditLogRepo IRepositoryMongoDB.ICommentEditLogsRepository) *DeleteCommentUsecase {
	return &DeleteCommentUsecase{
		eventBus:           eventBus,
		commentRepo:        commentRepo,
		commentEditLogRepo: commentEditLogRepo,
	}
}

func (d *DeleteCommentUsecase) Execute(ctx context.Context, req *req.DeleteCommentRequest) (*response.Response, error) {
	err := d.commentRepo.DeleteComment(ctx, req.CommentID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(400)), err
	}
	err = d.commentEditLogRepo.DeleteEditLogsByTargetID(ctx, req.CommentID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(400)), err
	}
	return nil, nil
}
