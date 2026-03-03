package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
)

type ICreatePageUsecase interface {
	Execute(ctx context.Context, req *req.CreateAdAccountRequest) (*res.FailedPageResponse, error)
}
type IUpdatePageUsecase interface {
	Execute(ctx context.Context, req *req.UpdateAdAccountRequest) (*res.FailedPageResponse, error)
}
type IDeletePageUsecase interface {
	Execute(ctx context.Context, req *req.DeletePageRequest) (*res.FailedPageResponse, error)
}
type IStatsPageUsecase interface {
	Execute(ctx context.Context) error
}
type ICreatePageRoleUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdatePageRoleUsecase interface {
	Execute(ctx context.Context) error
}
type IDeletePageRoleUsecase interface {
	Execute(ctx context.Context) error
}
type ICreatePageFollowerUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdatePageFollowerUsecase interface {
	Execute(ctx context.Context) error
}
type IDeletePageFollowerUsecase interface {
	Execute(ctx context.Context) error
}
type IMetricUseCase interface {
	Execute(ctx context.Context) error
}
type ICreateAdaccountUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdateAdaccountUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdateBalanceUsecase interface {
	Execute(ctx context.Context) error
}
type ICreateAdCampainUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdateAdCampainUsecase interface {
	Execute(ctx context.Context) error
}
type IDeleteAdCampainUsecase interface {
	Execute(ctx context.Context) error
}
type ICreateadsUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdateadsUsecase interface {
	Execute(ctx context.Context) error
}
type IDeleteadsUsecase interface {
	Execute(ctx context.Context) error
}
type IRankingAdCampainUsecase interface {
	Execute(ctx context.Context) error
}
type usecase struct {
}

func NewUsecase() *usecase {
	return &usecase{}
}
