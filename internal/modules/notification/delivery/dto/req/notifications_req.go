package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// NotificationReq chứa 100% các trường từ Entity.
// Best Practice: Sử dụng string cho các trường UUID để API nhận JSON chuẩn, sau đó parse ở Mapper.
type NotificationReq struct {
	UserID         string                       `json:"user_id" validate:"required,uuid"`
	CreatedAt      time.Time                    `json:"created_at"`
	NotificationID string                       `json:"notification_id" validate:"omitempty,uuid"` // Thường tự sinh, nhưng vẫn để để đủ 100%
	Type           sharedEnums.NotificationType `json:"type" validate:"required"`
	ActorID        string                       `json:"actor_id" validate:"required,uuid"`
	ActorName      string                       `json:"actor_name" validate:"required"`
	ActorAvatar    string                       `json:"actor_avatar"`
	TargetID       string                       `json:"target_id" validate:"required"`
	TargetPreview  string                       `json:"target_preview"`
	IsRead         bool                         `json:"is_read"`
	IsClicked      bool                         `json:"is_clicked"`
	GroupKey       string                       `json:"group_key"`
}
