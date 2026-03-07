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

type CreatePageRoleUsecase struct {
	pageRepo     IRepositoryMongodb.IPagesRepository
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
	events       events.EventBus
}

func NewCreatePageRoleUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRoleRepo IRepositoryMongodb.IPageRolesRepository, events events.EventBus) *CreatePageRoleUsecase {
	return &CreatePageRoleUsecase{
		pageRepo:     pageRepo,
		pageRoleRepo: pageRoleRepo,
		events:       events,
	}
}

func (u *CreatePageRoleUsecase) Execute(ctx context.Context, req *req.CreatePageRoleRequest) (*res.FailedPageRoleResponse, error) {
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
	entity := mapper.ToEntityPageRole(req.PageRoleReq)
	err = u.pageRoleRepo.CreatePageRole(ctx, entity)
	if err != nil {
		return &res.FailedPageRoleResponse{
			PageID:       req.PageID,
			UserID:       req.UserID,
			ErrorMessage: "Failed to create page role",
		}, nil
	}

	return &res.FailedPageRoleResponse{
		PageID:       req.PageID,
		UserID:       req.UserID,
		ErrorMessage: "Failed to create page role",
	}, nil
}
