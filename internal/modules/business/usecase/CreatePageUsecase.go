package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePageUsecase struct {
	pageRepo IRepositoryMongodb.IPagesRepository
	pageRole IRepositoryMongodb.IPageRolesRepository
}

func NewCreatePageUsecase() *CreatePageUsecase {
	return &CreatePageUsecase{}
}

func (u *CreatePageUsecase) Execute(ctx context.Context, req *req.CreateAdAccountRequest) (*res.FailedPageResponse, error) {
	entitypage := mapper.ToEntityPages(req.PageReq)
	err := u.pageRepo.CreatePage(ctx, entitypage)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       entitypage.ID.Hex(),
			UserActionID: req.CreatorUserID,
			ErrorMessage: err.Error(),
		}, err
	}

	entityRole := &entity.PageRole{
		ID:                primitive.NewObjectID(),
		PageID:            entitypage.ID,
		UserID:            entitypage.CreatorUserID,
		Role:              sharedEnums.RoleTypeAdmin,
		CustomPermissions: []string{"owner"},
		CreatedAt:         entitypage.CreatedAt,
		UpdatedAt:         entitypage.CreatedAt,
		AssignedBy:        entitypage.CreatorUserID,
	}
	err = u.pageRole.CreatePageRole(ctx, entityRole)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       entitypage.ID.Hex(),
			UserActionID: req.CreatorUserID,
			ErrorMessage: "Page created but failed to assign role to creator: " + err.Error(),
		}, err
	}
	return &res.FailedPageResponse{
		PageID:       entitypage.ID.Hex(),
		UserActionID: req.CreatorUserID,
		ErrorMessage: "Page created successfully",
	}, nil
}
