package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNotificationTemplates = "notification_templates"
)

type NotificationTemplate struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. CLASSIFICATION
	// Loại thông báo (VD: POST_LIKE)
	// Index: Unique { type: 1 } -> Mỗi loại chỉ có 1 template cấu hình
	Type sharedEnums.NotificationType `bson:"type" json:"type"`

	// 2. CONTENT (Multi-language)
	// Map key là mã ngôn ngữ ("vi", "en"), value là nội dung template
	// VD: "vi": "<strong>{{actor_name}}</strong> đã thích bài viết..."
	Template map[string]string `bson:"template" json:"template"`

	// 3. ASSETS & ACTIONS
	IconURL string `bson:"icon_url" json:"icon_url"`

	// Deeplink chứa placeholder
	// VD: "/post/{{target_id}}"
	ActionLink string `bson:"action_link" json:"action_link"`

	// 4. METADATA
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (NotificationTemplate) CollectionName() string {
	return CollectionNotificationTemplates
}
