package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotifGroupInvite struct {
}

func NewNotifGroupInvite() *NotifGroupInvite {
	return &NotifGroupInvite{}
}
func (n *NotifGroupInvite) Execute(ctx context.Context) error {
	return nil
}
func (n *NotifGroupInvite) GetType() sharedEnums.NotificationType {
	return sharedEnums.NotifGroupInvite
}
