package usecase

import (
	"context"

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

func (uc *HideOrUnhidePostUseCase) Execute(ctx context.Context, req *req.HideOrUnhidePostRequest) error {
	entity, err := uc.postRepo.GetPostByID(ctx, req.Post.ID)
	if err != nil {
		return err
	}
	if entity == nil {
		return nil
	}
	mapper.UpdateToEntityPost(req.Post, entity)
	_, err = uc.postRepo.UpdatePost(ctx, entity)
	if err != nil {
		return err
	}
	dataedit, err := uc.editRepo.GetByTargetID(ctx, req.Post.ID)
	if err != nil {
		return err
	}
	if dataedit != nil {
		mapper.UpdatePostEntityEditLogEntity(req.PostEdit, dataedit)
		err = uc.editRepo.UpdatePostEditLog(ctx, dataedit)
		if err != nil {
			return err
		}
	}
	if dataedit == nil {
		logEntity := mapper.ToPostEntityEditLogEntity(req.PostEdit)
		logEntity.TargetID = entity.ID
		err = uc.editRepo.CreatePostEditLog(ctx, logEntity)
		if err != nil {
			return err
		}
	}
	return nil
}
