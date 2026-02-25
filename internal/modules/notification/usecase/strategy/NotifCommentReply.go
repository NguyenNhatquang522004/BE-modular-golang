package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotifCommentReply struct {
}

func NewNotifCommentReply() *NotifCommentReply {
	return &NotifCommentReply{}
}
func (n *NotifCommentReply) Execute(ctx context.Context) error {
	// Implement the logic for executing the notification strategy
	return nil
}

func (n *NotifCommentReply) GetType() sharedEnums.NotificationType {
	// Return the type of notification
	return sharedEnums.NotifCommentReply
}
