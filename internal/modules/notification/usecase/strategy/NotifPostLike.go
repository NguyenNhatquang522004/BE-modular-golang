package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotiPostLike struct {
}

func NewNotifPostLike() *NotiPostLike {
	return &NotiPostLike{}
}
func (n *NotiPostLike) Execute(ctx context.Context) error {
	// Implement the logic for executing the notification strategy
	return nil
}
func (n *NotiPostLike) GetType() sharedEnums.NotificationType {
	return sharedEnums.NotifPostLike
}
