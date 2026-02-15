package usecase

type IFollowUseCase interface {
}

type IFriendshipUseCase interface {
}

type IBlockUseCase interface {
}

type IProfileUseCase interface {
}

type IAdminSocialUseCase interface {
}

type UseCase struct {
	FollowUseCase     IFollowUseCase
	FriendshipUseCase IFriendshipUseCase
	BlockUseCase      IBlockUseCase
	ProfileUseCase    IProfileUseCase
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
		FollowUseCase:     followUseCase,
		FriendshipUseCase: friendshipUseCase,
		BlockUseCase:      blockUseCase,
		ProfileUseCase:    profileUseCase,
		AdminSocialUseCase: adminSocialUseCase,
	}
}
