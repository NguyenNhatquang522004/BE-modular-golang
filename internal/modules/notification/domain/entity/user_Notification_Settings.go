package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
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

	Name        string `bson:"name" json:"name"` // Tên người dùng (để hiển thị trong noti nếu cần)
	Avatar      string `bson:"avatar" json:"avatar"`
	DateOfBirth string `bson:"date_of_birth" json:"date_of_birth"`

	// Cấu hình
	Settings *GeneralSettings `bson:"settings" json:"settings"`

	// Danh sách FCM Token (Array)
	// Index: Multikey { "fcm_tokens.token": 1 } -> Để tìm và xóa token chết (invalid)da
	FCMTokens []*FCMToken `bson:"fcm_tokens,omitempty" json:"fcm_tokens,omitempty"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time  `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- FCM TOKEN (Mobile Push) ---
type FCMToken struct {
	Token     string    `bson:"token" json:"token"`         // Token firebase gửi về
	DeviceID  string    `bson:"device_id" json:"device_id"` // UUID của thiết bị (để replace token cũ của máy đó)
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// --- GENERAL SETTINGS ---
type GeneralSettings struct {
	PushEnabled      bool                       `bson:"push_enabled" json:"push_enabled"` // Master switch (Tắt tất cả push)
	EmailFrequency   sharedEnums.EmailFrequency `bson:"email_frequency" json:"email_frequency"`
	PushInteractions bool                       `bson:"push_interactions" json:"push_interactions"`
	PushFriends      bool                       `bson:"push_friends" json:"push_friends"`
	PushGroups       bool                       `bson:"push_groups" json:"push_groups"`
	PushEvents       bool                       `bson:"push_events" json:"push_events"`
	PushBirthdays    bool                       `bson:"push_birthdays" json:"push_birthdays"`
}

func (UserNotificationSetting) CollectionName() string {
	return CollectionUserNotificationSettings
}
