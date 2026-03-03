package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
)

type UpdatePageUsecase struct {

	pageRepo IRepositoryMongodb.IPagesRepository
	pageRole IRepositoryMongodb.IPageRolesRepository
}

func NewUpdatePageUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRole IRepositoryMongodb.IPageRolesRepository) *UpdatePageUsecase {
	return &UpdatePageUsecase{
		pageRepo: pageRepo,
		pageRole: pageRole,
	}
}
func (u *UpdatePageUsecase) Execute(ctx context.Context, req *req.UpdateAdAccountRequest) (*res.FailedPageResponse, error) {
	datapage, err := u.pageRepo.GetPageByID(ctx, req.ID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.ID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to retrieve page: " + err.Error(),
		}, err
	}
	if datapage == nil {
		return &res.FailedPageResponse{
			PageID:       req.ID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Page not found",
		}, nil
	}
	dataUserAction, err := u.pageRole.GetPageRolesByPageIDAndUserID(ctx, datapage.ID.Hex(), req.UserActionID)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.ID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to retrieve user role: " + err.Error(),
		}, err
	}
	if dataUserAction == nil {
		return &res.FailedPageResponse{
			PageID:       req.ID,
			UserActionID: req.UserActionID,
			ErrorMessage: "User has no role on this page",
		}, nil
	}
	if datapage.CreatorUserID != req.UserActionID {
		if dataUserAction.Role != sharedEnums.RoleTypeAdmin && dataUserAction.Role != sharedEnums.RoleTypeModerator {
			return &res.FailedPageResponse{
				PageID:       req.ID,
				UserActionID: req.UserActionID,
				ErrorMessage: "User does not have permission to update this page",
			}, nil
		}
	}

	mapper.UpdateToEntityPages(req.PageReq, datapage)
	err = u.pageRepo.UpdatePage(ctx, datapage)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       req.ID,
			UserActionID: req.UserActionID,
			ErrorMessage: "Failed to update page: " + err.Error(),
		}, err
	}

	return &res.FailedPageResponse{
		PageID:       req.ID,
		UserActionID: req.UserActionID,
		ErrorMessage: "Page updated successfully",
	}, nil
}
