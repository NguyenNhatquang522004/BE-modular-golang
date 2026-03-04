package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type IFollowUseCase interface {
	CreateFollowUserUseCase(ctx context.Context, req *req.FollowCreateRequest) (*response.Response, error)
	DeleteSoftFollowUserUseCase(ctx context.Context, follower *req.FollowDeleteSoftRequest) (*response.Response, error)
	DeleteHardFollowUserUseCase(ctx context.Context, follower *req.FollowDeleteHardRequest) (*response.Response, error)
	UpdatateMuteFollowUserUseCase(ctx context.Context, follower *req.FollowUpdateMuteRequest) (*response.Response, error)
	PaginationFollowersUseCase(ctx context.Context, req *req.FollowPaginationRequest) (*response.Response, error)
	PaginationFollowedsUseCase(ctx context.Context, req *req.FollowPaginationRequest) (*response.Response, error)
}

type IFriendshipUseCase interface {
	HandleFriendShipUseCase(ctx context.Context, req *req.FriendShipUseCaseRequest) (*response.Response, error)
	CreateFriendshipUseCase(ctx context.Context, req *req.CreateFriendshipRequest) (*response.Response, error)
	PanigationAcceptedFriendshipUseCase(ctx context.Context, req *req.PaginationFriendshipRequest) (*response.Response, error)
	PanigationPendingFriendshipUseCase(ctx context.Context, req *req.PaginationFriendshipRequest) (*response.Response, error)
}

type IBlockUseCase interface {
	UseCaseBlockUser(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error)
	GetPaginationTypeBlockUseCase(ctx context.Context, req *req.BlockPaginationTypeBlockRequest) (*response.Response, error)
	IsBlockedUseCase(ctx context.Context, req *req.BlockIsBlockedRequest) (*response.Response, error)
}

type IProfileUseCase interface {
	CreateProfileUseCase(ctx context.Context, req *req.CreateAndUpdateProfileRequest) (*response.Response, error)
	GetProfileByIDUseCase(ctx context.Context, req *req.ProfileIDRequest) (*response.Response, error)
	UpdateProfileUseCase(ctx context.Context, req *req.CreateAndUpdateProfileRequest) (*response.Response, error)
	GetProfileByUserIDUseCase(ctx context.Context, req *req.ProfileIDRequest) (*response.Response, error)
}

type IAdminSocialUseCase interface {
}

type UseCase struct {
	FollowUseCase      IFollowUseCase
	FriendshipUseCase  IFriendshipUseCase
	BlockUseCase       IBlockUseCase
	ProfileUseCase     IProfileUseCase
	AdminSocialUseCase IAdminSocialUseCase
}

func NewUseCase(
	followUseCase IFollowUseCase,
	friendshipUseCase IFriendshipUseCase,
	blockUseCase IBlockUseCase,
	profileUseCase IProfileUseCase,
	adminSocialUseCase IAdminSocialUseCase,
) *UseCase {
	return &UseCase{
		FollowUseCase:      followUseCase,
		FriendshipUseCase:  friendshipUseCase,
		BlockUseCase:       blockUseCase,
		ProfileUseCase:     profileUseCase,
		AdminSocialUseCase: adminSocialUseCase,
	}
}
