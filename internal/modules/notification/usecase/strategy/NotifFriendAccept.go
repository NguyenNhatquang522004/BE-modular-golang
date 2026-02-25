package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotifFriendAccept struct {
}

func NewNotifFriendAccept() *NotifFriendAccept {
	return &NotifFriendAccept{}
}
func (n *NotifFriendAccept) Execute(ctx context.Context) error {
	// Implement the logic for executing the notification strategy
	return nil
}
func (n *NotifFriendAccept) GetType() sharedEnums.NotificationType {
	// Return the type of notification
	return sharedEnums.NotifFriendAccept
}
