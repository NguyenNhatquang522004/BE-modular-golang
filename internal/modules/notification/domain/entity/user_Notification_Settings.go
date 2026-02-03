package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUserNotificationSettings = "user_notification_settings"
)

type UserNotificationSetting struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// User ID (Postgres UUID -> String)
	// Index: Unique { user_id: 1 } -> Mỗi user chỉ có 1 bản ghi settings
	UserID string `bson:"user_id" json:"user_id"`

	// Cấu hình
	Settings GeneralSettings `bson:"settings" json:"settings"`

	// Danh sách FCM Token (Array)
	// Index: Multikey { "fcm_tokens.token": 1 } -> Để tìm và xóa token chết (invalid)da
	FCMTokens []FCMToken `bson:"fcm_tokens,omitempty" json:"fcm_tokens,omitempty"`
}

// --- FCM TOKEN (Mobile Push) ---
type FCMToken struct {
	Token     string    `bson:"token" json:"token"`         // Token firebase gửi về
	DeviceID  string    `bson:"device_id" json:"device_id"` // UUID của thiết bị (để replace token cũ của máy đó)
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// --- CHANNEL TOGGLE (Cấu trúc tái sử dụng) ---
// Định nghĩa: Với sự kiện X, user muốn nhận qua kênh nào?
type ChannelToggle struct {
	Push  bool `bson:"push" json:"push"`
	Email bool `bson:"email" json:"email"`
}

// --- DETAIL SETTINGS ---
// Struct hóa các key thay vì dùng Map để đảm bảo Type Safety.
// Nếu sau này thêm 'on_friend_request', chỉ cần thêm field vào đây.
type NotificationDetailSettings struct {
	OnPostLike ChannelToggle `bson:"on_post_like" json:"on_post_like"`
	OnComment  ChannelToggle `bson:"on_comment" json:"on_comment"`
	OnBirthday ChannelToggle `bson:"on_birthday" json:"on_birthday"`

	// Mở rộng thêm:
	OnFriendRequest ChannelToggle `bson:"on_friend_request" json:"on_friend_request"`
	OnSystemAlert   ChannelToggle `bson:"on_system_alert" json:"on_system_alert"`
}

// --- GENERAL SETTINGS ---
type GeneralSettings struct {
	PushEnabled    bool                `bson:"push_enabled" json:"push_enabled"` // Master switch (Tắt tất cả push)
	EmailFrequency enum.EmailFrequency `bson:"email_frequency" json:"email_frequency"`

	// Chi tiết từng loại
	Details NotificationDetailSettings `bson:"details" json:"details"`
}

func (UserNotificationSetting) CollectionName() string {
	return CollectionUserNotificationSettings
}
