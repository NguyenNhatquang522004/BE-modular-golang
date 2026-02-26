package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
)

type IPublishPostUseCase interface {
	Execute(ctx context.Context, req *req.PublishPostRequest) (*response.Response, error)
}
type ISharePostUseCase interface {
	Execute(ctx context.Context, req *req.SharePostRequest) (*response.Response, error)
}
type IDeleteSharePostUseCase interface {
	Execute(ctx context.Context, req *req.SharePostRequest) (*response.Response, error)
}
type IDeleteSharePostToGroupUseCase interface { // IGNORE --- đợi làm tới community
	Execute(ctx context.Context, req *req.SharePostToGroupRequest) (*response.Response, error)
}
type ISharePostToGroupUseCase interface { // IGNORE --- đợi làm tới community
	Execute(ctx context.Context, req *req.SharePostToGroupRequest) (*response.Response, error)
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
type IGetEnumUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}

type UsecaseContent struct {
	pulishPostUseCase       IPublishPostUseCase
	editPostUseCase         IEditPostUseCase
	hideOrUnhidePostUseCase IHideOrUnhidePostUseCase
	deletePostUseCase       IDeletePostUseCase
	getPostByUserIDUseCase  IGetPostByUserIDUseCase
	enumUseCase             IGetEnumUseCase
}

func NewUsecaseContent(publishPostUC IPublishPostUseCase, editPostUC IEditPostUseCase, hideOrUnhidePostUC IHideOrUnhidePostUseCase, deletePostUC IDeletePostUseCase, getPostByUserIDUC IGetPostByUserIDUseCase, enumUC IGetEnumUseCase) *UsecaseContent {
	return &UsecaseContent{
		pulishPostUseCase:       publishPostUC,
		editPostUseCase:         editPostUC,
		hideOrUnhidePostUseCase: hideOrUnhidePostUC,
		deletePostUseCase:       deletePostUC,
		getPostByUserIDUseCase:  getPostByUserIDUC,
		enumUseCase:             enumUC,
	}
}
