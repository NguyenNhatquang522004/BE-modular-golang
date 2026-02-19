package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
)

type IFollowUseCase interface {
	CreateFollowUserUseCase(req *req.FollowCreateRequest) (*response.Response, error)
	DeleteSoftFollowUserUseCase(follower *req.FollowDeleteSoftRequest) (*response.Response, error)
	DeleteHardFollowUserUseCase(follower *req.FollowDeleteHardRequest) (*response.Response, error)
	UpdatateMuteFollowUserUseCase(follower *req.FollowUpdateMuteRequest) (*response.Response, error)
	PaginationFollowersUseCase(req *req.FollowPaginationRequest) (*response.Response, error)
	PaginationFollowedsUseCase(req *req.FollowPaginationRequest) (*response.Response, error)
}

type IFriendshipUseCase interface {
	HandleFriendShipUseCase(req *req.FriendShipUseCaseRequest) (*response.Response, error)
	CreateFriendshipUseCase(req *req.CreateFriendshipRequest) (*response.Response, error)
	PanigationAcceptedFriendshipUseCase(req *req.PaginationFriendshipRequest) (*response.Response, error)
	PanigationPendingFriendshipUseCase(req *req.PaginationFriendshipRequest) (*response.Response, error)
}

type IBlockUseCase interface {
	UseCaseBlockUser(req *req.BlockCreateRequest) (*response.Response, error)
	GetPaginationTypeBlockUseCase(req *req.BlockPaginationTypeBlockRequest) (*response.Response, error)
	IsBlockedUseCase(req *req.BlockIsBlockedRequest) (*response.Response, error)
}

type IProfileUseCase interface {
	CreateProfileUseCase(req *req.CreateAndUpdateProfileRequest) (*response.Response, error)
	GetProfileByIDUseCase(req *req.ProfileIDRequest) (*response.Response, error)
	UpdateProfileUseCase(req *req.CreateAndUpdateProfileRequest) (*response.Response, error)
	GetProfileByUserIDUseCase(req *req.ProfileIDRequest) (*response.Response, error)
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
