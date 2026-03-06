package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type IFollowUseCase interface {
	Execute(ctx context.Context, req *req.FollowerRequest) error
}
type IFriendshipUseCase interface {
	Execute(ctx context.Context, req *req.FriendShipRequest) error
}
type IBlockUseCase interface {
	Execute(ctx context.Context, req *req.BlockRequest) error
}
type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}
