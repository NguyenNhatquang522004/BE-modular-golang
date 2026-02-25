package notificationEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"

type NotificationEvent struct {
	TypeNotification []*sharedEnums.NotificationType
	SendUser         []string
	SendGroup        []string
	SendPost         []string
}
