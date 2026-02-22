package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
)

// --- FCM TOKEN RES ---
type FCMTokenRes struct {
	Token     string    `json:"token"`
	DeviceID  string    `json:"device_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- GENERAL SETTINGS RES ---
type GeneralSettingsRes struct {
	PushEnabled      bool                `json:"push_enabled"`
	EmailFrequency   *enum.EmailFrequency `json:"email_frequency"`
	PushInteractions bool                `json:"push_interactions"`
	PushFriends      bool                `json:"push_friends"`
	PushGroups       bool                `json:"push_groups"`
	PushEvents       bool                `json:"push_events"`
	PushBirthdays    bool                `json:"push_birthdays"`
}

// --- MAIN RESPONSE ---
type UserNotificationSettingRes struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	Settings  *GeneralSettingsRes `json:"settings"`
	FCMTokens []*FCMTokenRes      `json:"fcm_tokens,omitempty"`
}