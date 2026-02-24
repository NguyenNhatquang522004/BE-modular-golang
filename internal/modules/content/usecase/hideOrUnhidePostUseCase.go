package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type HideOrUnhidePostUseCase struct {
	postRepo IRepositoryMongodb.IPostRepository
	editRepo IRepositoryMongodb.IPostEditLogsRepository
}

func NewHideOrUnhidePostUseCase(postRepo IRepositoryMongodb.IPostRepository, postMediaRepo IRepositoryMongodb.IPostMediaRepository, postSetting IRepositoryMongodb.IPostSettingRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postInsightRepo IRepositoryCassandra.IPostInsights) *HideOrUnhidePostUseCase {
	return &HideOrUnhidePostUseCase{
		postRepo: postRepo,
	}
}

func (uc *HideOrUnhidePostUseCase) Execute(ctx context.Context, req *req.HideOrUnhidePostRequest) (*response.Response, error) {
	entity, err := uc.postRepo.GetPostByID(ctx, req.Post.ID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithMessage("Post not found"),
			response.WithStatus(http.StatusNotFound),
		), nil
	}
	mapper.UpdateToEntityPost(req.Post, entity)
	_, err = uc.postRepo.UpdatePost(ctx, entity)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage("Failed to update post"), response.WithStatus(http.StatusBadRequest)), err
	}
	dataedit, err := uc.editRepo.GetByTargetID(ctx, req.Post.ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage("Failed to update post"), response.WithStatus(http.StatusBadRequest)), err
	}
	if dataedit != nil {
		mapper.UpdatePostEntityEditLogEntity(req.PostEdit, dataedit)
		err = uc.editRepo.UpdatePostEditLog(ctx, dataedit)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage("Failed to update post"), response.WithStatus(http.StatusBadRequest)), err
		}
	}
	if dataedit == nil {
		logEntity := mapper.ToPostEntityEditLogEntity(req.PostEdit)
		logEntity.TargetID = entity.ID
		err = uc.editRepo.CreatePostEditLog(ctx, logEntity)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage("Failed to update post"), response.WithStatus(http.StatusBadRequest)), err
		}
	}
	return response.NewResponse(
		response.WithData(""),
		response.WithMessage("Post updated successfully"),
		response.WithStatus(http.StatusOK),
	), nil
}
