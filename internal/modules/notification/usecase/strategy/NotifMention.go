package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotifMention struct {
}

func NewNotifMention() *NotifMention {
	return &NotifMention{}
}
func (n *NotifMention) Execute(ctx context.Context)  error{
	// Implement the logic for executing the notification strategy
	return nil
}
func (n *NotifMention) GetType() sharedEnums.NotificationType {
	// Return the type of notification
	return sharedEnums.NotifMention
}
