package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// NotificationRes dùng để trả dữ liệu cho Client.
// UUID được chuyển thành string để frontend dễ dàng xử lý.
type NotificationRes struct {
	UserID         string                       `json:"user_id"`
	CreatedAt      time.Time                    `json:"created_at"`
	NotificationID string                       `json:"notification_id"`
	Type           sharedEnums.NotificationType `json:"type"`
	ActorID        string                       `json:"actor_id"`
	ActorName      string                       `json:"actor_name"`
	ActorAvatar    string                       `json:"actor_avatar"`
	TargetID       string                       `json:"target_id"`
	TargetPreview  string                       `json:"target_preview"`
	IsRead         bool                         `json:"is_read"`
	IsClicked      bool                         `json:"is_clicked"`
	GroupKey       string                       `json:"group_key"`
}
