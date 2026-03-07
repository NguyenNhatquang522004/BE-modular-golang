package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type UpdatePageRoleUsecase struct {
	pageRepo     IRepositoryMongodb.IPagesRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	events       events.EventBus
}

func NewUpdatePageRoleUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, events events.EventBus) *UpdatePageRoleUsecase {
	return &UpdatePageRoleUsecase{
		pageRepo:     pageRepo,
		pageRoleRepo: pageRoleRepo,
		events:       events,
	}
}

func (u *UpdatePageRoleUsecase) Execute(ctx context.Context, req *req.UpdatePageRoleRequest) (*res.FailedPageRoleResponse, error) {
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
	dataUserAction, err := u.pageRoleRepo.GetPageRoleByPageIDAndUserID(ctx, datapage.ID.Hex(), req.UserActionID)
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
	datauser, err := u.pageRoleRepo.GetPageRoleByPageIDAndUserID(ctx, datapage.ID.Hex(), req.UserID)
	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to get user role",
		}, nil
	}
	if datauser == nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "User does not have a role on this page",
		}, nil
	}
	if dataUserAction.Role == datauser.Role {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "User Not Permission To Update Role",
		}, nil
	}
	mapper.UpdateToEntityPageRole(req.PageRoleReq, datauser)
	err = u.pageRoleRepo.UpdatePageRole(ctx, datauser)
	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to update page role",
		}, nil
	}
	return &res.FailedPageRoleResponse{
		PageID:       req.PageID,
		UserID:       req.UserID,
		ErrorMessage: "Failed to update page role",
	}, nil
}
