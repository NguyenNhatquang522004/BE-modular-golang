package usecase

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePageUsecase struct {
	pageRepo      IRepositoryMongodb.IPagesRepository
	pageRole      IRepositoryMongodb.IPageRolesRepository
	pageDailyRepo IRepositoryCassandra.IPageDailyMetricsRepository
}

func NewCreatePageUsecase(pageRepo IRepositoryMongodb.IPagesRepository, pageRole IRepositoryMongodb.IPageRolesRepository, pageDailyRepo IRepositoryCassandra.IPageDailyMetricsRepository) *CreatePageUsecase {
	return &CreatePageUsecase{
		pageRepo:      pageRepo,
		pageRole:      pageRole,
		pageDailyRepo: pageDailyRepo,
	}
}

func (u *CreatePageUsecase) Execute(ctx context.Context, req *req.CreatePageRequest) (*res.FailedPageResponse, error) {
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
	convertedPageIDCQL, err := gocql.ParseUUID(entitypage.ID.Hex())
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       entitypage.ID.Hex(),
			UserActionID: req.CreatorUserID,
			ErrorMessage: "Page created but failed to parse Page ID for daily metrics: " + err.Error(),
		}, err
	}

	entityMetric := &entity.PageDailyMetric{
		ID:               convertedPageIDCQL,
		MetricDate:       entitypage.CreatedAt.Truncate(24 * time.Hour), // Lưu ý: MetricDate chỉ lưu ngày, không lưu giờ
		ReachTotal:       0,
		ReachPaid:        0,
		ReachOrganic:     0,
		ImpressionsTotal: 0,
		NewFollowers:     0,
		Unfollows:        0,
		ProfileViews:     0,
		WebsiteClicks:    0,
		CTAClicks:        0,
	}
	err = u.pageDailyRepo.CreatePageDailyMetric(ctx, entityMetric)
	if err != nil {
		return &res.FailedPageResponse{
			PageID:       entitypage.ID.Hex(),
			UserActionID: req.CreatorUserID,
			ErrorMessage: "Page created but failed to create initial daily metrics: " + err.Error(),
		}, err
	}
	return &res.FailedPageResponse{
		PageID:       entitypage.ID.Hex(),
		UserActionID: req.CreatorUserID,
		ErrorMessage: "Page created successfully",
	}, nil
}
