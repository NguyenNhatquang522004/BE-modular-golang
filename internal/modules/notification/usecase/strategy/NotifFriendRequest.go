package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotifFriendRequest struct {
}

func NewNotifFriendRequest() *NotifFriendRequest {
	return &NotifFriendRequest{}
}
func (n *NotifFriendRequest) Execute(ctx context.Context)  error{
	// Implement the logic for executing the notification strategy 
	return nil
}
func (n *NotifFriendRequest) GetType() sharedEnums.NotificationType {
	// Return the type of notification
	return sharedEnums.NotifFriendRequest
}
