package usecase

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"

type FriendshipUseCase struct {
	events events.EventBus
}

func NewFriendshipUseCase(events events.EventBus) *FriendshipUseCase {

	return &FriendshipUseCase{
		events: events,
	}
}
