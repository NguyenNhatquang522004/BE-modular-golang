package strategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotiFSystemAlert struct {
}

func NewNotifSystemAlert() *NotiFSystemAlert {
	return &NotiFSystemAlert{}
}
func (n *NotiFSystemAlert) Execute(ctx context.Context)  error {
	// Implement the logic for executing the notification strategy
	return nil
}
func (n *NotiFSystemAlert) GetType() sharedEnums.NotificationType {
	return sharedEnums.NotifSystemAlert
}
