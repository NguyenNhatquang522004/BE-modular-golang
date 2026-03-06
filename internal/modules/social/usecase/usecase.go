package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type IFollowUseCase interface {
	Execute(ctx context.Context, req *req.FollowerRequest) error
}
type IFriendshipUseCase interface {
}
type IBlockUseCase interface {
}
type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}
