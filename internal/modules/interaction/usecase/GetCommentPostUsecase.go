package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type GetCommentPostUsecase struct {
	commentRepo IRepositoryMongoDB.ICommentRepository
}

func NewGetCommentPostUsecase(commentRepo IRepositoryMongoDB.ICommentRepository) *GetCommentPostUsecase {
	return &GetCommentPostUsecase{
		commentRepo: commentRepo,
	}
}

func (u *GetCommentPostUsecase) Execute(ctx context.Context, req *req.GetCommentPostRequest) (*response.Response, error) {
	// Implement the logic for getting comments of a post here
	comments, err := u.commentRepo.PaginationComments(ctx, req.TargetID, req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	return response.NewResponse(response.WithData(comments), response.WithMessage("Comments retrieved successfully"), response.WithStatus(http.StatusOK)), nil
}
