package usecase

import "context"

type ICreatePageUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdatePageUsecase interface {
	Execute(ctx context.Context) error
}
type IDeletePageUsecase interface {
	Execute(ctx context.Context) error
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
type IMetricUseCase interface{
	Execute(ctx context.Context) error
}


type usecase struct {
}

func NewUsecase() *usecase {
	return &usecase{}
}
