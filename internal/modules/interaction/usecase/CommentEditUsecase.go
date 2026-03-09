package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type CommentEditUsecase struct {
	CommentEditLogRepo IRepositoryMongoDB.IEditLogsRepository
	commentRepo        IRepositoryMongoDB.ICommentRepository
	pool               IRepositoryShare.IWorkerPool
}

func NewCommentEditUsecase(commentEditLogRepo IRepositoryMongoDB.IEditLogsRepository, commentRepo IRepositoryMongoDB.ICommentRepository, pool IRepositoryShare.IWorkerPool) *CommentEditUsecase {
	return &CommentEditUsecase{
		CommentEditLogRepo: commentEditLogRepo,
		commentRepo:        commentRepo,
		pool:               pool,
	}
}

func (u *CommentEditUsecase) Execute(ctx context.Context, req *req.CommentEditUseRequest) (*response.Response, error) {
	// Implement the logic for editing a comment here
	data, err := u.CommentEditLogRepo.GetLatestEditLogByTargetID(ctx, req.TargetID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	req.CommentEditLogReq.Version = data.Version + 1
	entity, err := mapper.ToEntityCommentEditLogs(req.CommentEditLogReq)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	err = u.CommentEditLogRepo.CreateEditLog(ctx, entity)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("Comment edited successfully"), response.WithStatus(http.StatusOK)), nil
}
