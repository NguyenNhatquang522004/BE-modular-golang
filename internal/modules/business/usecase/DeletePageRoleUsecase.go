package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type DeletePageRoleUsecase struct {
	pageRepo     IRepositoryMongodb.IPagesRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	events       events.EventBus
}

func NewDeletePageRoleUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, events events.EventBus) *DeletePageRoleUsecase {
	return &DeletePageRoleUsecase{
		pageRepo:     pageRepo,
		pageRoleRepo: pageRoleRepo,
		events:       events,
	}
}

func (u *DeletePageRoleUsecase) Execute(ctx context.Context, req *req.DeletePageRoleRequest) (*res.FailedPageRoleResponse, error) {
	datapage, err := u.pageRepo.GetPageByID(ctx, req.PageID)
	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Page not found",
		}, nil
	}
	if datapage == nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Page not found",
		}, nil
	}
	dataUserAction, err := u.pageRoleRepo.GetPageRolesByPageIDAndUserID(ctx, datapage.ID.Hex(), req.UserActionID)
	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to get user action",
		}, nil
	}
	if dataUserAction.UserID != datapage.CreatorUserID {
		if dataUserAction.Role != sharedEnums.RoleTypeAdmin {
			return &res.FailedPageRoleResponse{
				PageID:       req.PageID,
				UserID:       req.UserID,
				ErrorMessage: "User does not have permission to assign roles",
			}, nil
		}
	}
	err = u.pageRoleRepo.DeletePageRoleByUserIDandPageID(ctx, req.UserID, req.PageID)

	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to delete page role",
		}, nil
	}
	return &res.FailedPageRoleResponse{
		PageID:       req.PageID,
		UserID:       req.UserID,
		ErrorMessage: "",
	}, nil
}
