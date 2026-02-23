package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
)

type IPublishPostUseCase interface {
	Execute(ctx context.Context, req *req.PublishPostRequest) (*response.Response, error)
	ExecuteBulk(ctx context.Context, req []*req.PublishPostRequest) (*response.Response, error)
}

type IEditPostUseCase interface {
	Execute(ctx context.Context, req *req.EditPostRequest) (*response.Response, error)
}

type IHideOrUnhidePostUseCase interface {
	Execute(ctx context.Context, req *req.HideOrUnhidePostRequest) (*response.Response, error)
}

type IDeletePostUseCase interface {
	Execute(ctx context.Context, req *req.DeletePostRequest) (*response.Response, error)
}

type IGetPostByUserIDUseCase interface {
	Execute(ctx context.Context, req *req.GetPostByUserIDRequest) (*response.Response, error)
}

type IGetPostDetailUseCase interface {
	Execute(ctx context.Context, req *req.GetPostDetailRequest) (*response.Response, error)
}
type UsecaseContent struct {
	pulishPostUseCase       IPublishPostUseCase
	editPostUseCase         IEditPostUseCase
	hideOrUnhidePostUseCase IHideOrUnhidePostUseCase
	deletePostUseCase       IDeletePostUseCase
	getPostByUserIDUseCase  IGetPostByUserIDUseCase
	getPostDetailUseCase    IGetPostDetailUseCase
}

func NewUsecaseContent(publishPostUC IPublishPostUseCase, editPostUC IEditPostUseCase, hideOrUnhidePostUC IHideOrUnhidePostUseCase, deletePostUC IDeletePostUseCase, getPostByUserIDUC IGetPostByUserIDUseCase, getPostDetailUC IGetPostDetailUseCase) *UsecaseContent {
	return &UsecaseContent{
		pulishPostUseCase:       publishPostUC,
		editPostUseCase:         editPostUC,
		hideOrUnhidePostUseCase: hideOrUnhidePostUC,
		deletePostUseCase:       deletePostUC,
		getPostByUserIDUseCase:  getPostByUserIDUC,
		getPostDetailUseCase:    getPostDetailUC,
	}
}
