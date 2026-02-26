package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
)

type ICommentPostAndReplyUsecase interface {
	Execute(ctx context.Context, req *req.CommentPostAndReplyRequest) (*response.Response, error)
}

type ICommentEditUsecase interface {
	Execute(ctx context.Context, req *req.CommentEditUseRequest) (*response.Response, error)
}
type ICommentReactionUsecase interface {
	Execute(ctx context.Context, req *interactionEvent.CommentReactionPayload) (*response.Response, error)
}
type IReactionPostUsecase interface {
	Execute(ctx context.Context, req *interactionEvent.PostReactionPayload) (*response.Response, error)
}

type IDeleteCommentUsecase interface {
	Execute(ctx context.Context, req *req.DeleteCommentRequest) (*response.Response, error)
}
type IBookmarkPostUsecase interface {
	Execute(ctx context.Context, req *req.BookmarkPostRequest) (*response.Response, error)
}
type IDeleteBookmarkUsecase interface {
	Execute(ctx context.Context, req *req.DeleteBookmarkRequest) (*response.Response, error)
}
type IGetCommentPostUsecase interface {
	Execute(ctx context.Context, req *req.GetCommentPostRequest) (*response.Response, error)
}

type IGetCommentDetailUsecase interface {
	Execute(ctx context.Context, req *req.GetCommentDetailRequest) (*response.Response, error)
}
type IGetBookMarkPostUsecase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type Usecase struct {
	ICommentPostAndReplyUsecase ICommentPostAndReplyUsecase
}

func NewUsecase(commentPostAndReplyUsecase ICommentPostAndReplyUsecase) *Usecase {
	return &Usecase{
		ICommentPostAndReplyUsecase: commentPostAndReplyUsecase,
	}
}
