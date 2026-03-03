package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type DeletePageUsecase struct {
	pageRepo   IRepositoryMongodb.IPagesRepository
	pageRole   IRepositoryMongodb.IPageRolesRepository
	pageFollow IRepositoryMongodb.IPageFollowersRepository
}

func NewDeletePageUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRole IRepositoryMongodb.IPageRolesRepository, pageFollow IRepositoryMongodb.IPageFollowersRepository) *DeletePageUsecase {
	return &DeletePageUsecase{
		pageRepo:   pageRepo,
		pageRole:   pageRole,
		pageFollow: pageFollow,
	}
}

func (u *DeletePageUsecase) Execute(ctx context.Context, req *req.DeletePageRequest) (*res.FailedPageResponse, error) {

	return &res.FailedPageResponse{
		PageID:       req.PageID,
		UserActionID: req.UserActionID,
		ErrorMessage: "Delete page ",
	}, nil
}
