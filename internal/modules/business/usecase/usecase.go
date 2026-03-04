package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
)

type ICreatePageUsecase interface {
	Execute(ctx context.Context, req *req.CreatePageRequest) (*res.FailedPageResponse, error)
}
type IUpdatePageUsecase interface {
	Execute(ctx context.Context, req *req.UpdatePageRequest) (*res.FailedPageResponse, error)
}
type IDeletePageUsecase interface {
	Execute(ctx context.Context, req *req.DeletePageRequest) (*res.FailedPageResponse, error)
}
type IStatsPageUsecase interface {
	Execute(ctx context.Context, req *req.StatsPageRequest) (*res.FailedPageResponse, error)
}
type ICreatePageRoleUsecase interface {
	Execute(ctx context.Context, req *req.CreatePageRoleRequest) (*res.FailedPageRoleResponse, error)
}
type IUpdatePageRoleUsecase interface {
	Execute(ctx context.Context, req *req.UpdatePageRoleRequest) (*res.FailedPageRoleResponse, error)
}
type IDeletePageRoleUsecase interface {
	Execute(ctx context.Context, req *req.DeletePageRoleRequest) (*res.FailedPageRoleResponse, error)
}
type ICreatePageFollowerUsecase interface {
	Execute(ctx context.Context, req *req.CreatePageFollowerRequest) (*res.FailedPageFollowerResponse, error)
}
type IUpdatePageFollowerUsecase interface {
	Execute(ctx context.Context, req *req.CreatePageFollowerRequest) (*res.FailedPageFollowerResponse, error)
}
type IDeletePageFollowerUsecase interface {
	Execute(ctx context.Context, req *req.CreatePageFollowerRequest) (*res.FailedPageFollowerResponse, error)
}
type IMetricUseCase interface {
	Execute(ctx context.Context) error // làm sau
}
type ICreateAdaccountUsecase interface {
	Execute(ctx context.Context, req *req.CreateAdAccountRequest) (*res.FailedAdAccountResponse, error)
}
type IUpdateAdaccountUsecase interface {
	Execute(ctx context.Context, req *req.UpdateAdAccountRequest) (*res.FailedAdAccountResponse, error)
}
type IUpdateBalanceUsecase interface {
	Execute(ctx context.Context, req *req.UpdateBalanceRequest) (*res.FailedAdAccountResponse, error)
}
type ICreateAdCampainUsecase interface {
	Execute(ctx context.Context, req *req.CreateAdCampainRequest) (*res.FailedAdCampainResponse, error)
}
type IUpdateAdCampainUsecase interface {
	Execute(ctx context.Context, req *req.UpdateAdCampainRequest) (*res.FailedAdCampainResponse, error)
}
type IDeleteAdCampainUsecase interface {
	Execute(ctx context.Context, req *req.DeleteAdCampainRequest) (*res.FailedAdCampainResponse, error)
}
type ICreateadsUsecase interface {
	Execute(ctx context.Context, req *req.CreateAdsRequest) (*res.FailedAdResponse, error)
}
type IUpdateadsUsecase interface {
	Execute(ctx context.Context, req *req.UpdateAdsRequest) (*res.FailedAdResponse, error)
}
type IDeleteadsUsecase interface {
	Execute(ctx context.Context, req *req.DeleteAdsRequest) (*res.FailedAdResponse, error)
}
type IRankingAdCampainUsecase interface {
	Execute(ctx context.Context) error
}
type usecase struct {
}

func NewUsecase() *usecase {
	return &usecase{}
}
