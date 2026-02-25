package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
)

// --- FCM TOKEN REQ ---
type FCMTokenReq struct {
	Token     string    `json:"token" binding:"required"`
	DeviceID  string    `json:"device_id" binding:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- GENERAL SETTINGS REQ ---
type GeneralSettingsReq struct {
	PushEnabled      bool                 `json:"push_enabled"`
	EmailFrequency   *enum.EmailFrequency `json:"email_frequency" binding:"required"`
	PushInteractions bool                 `json:"push_interactions"`
	PushFriends      bool                 `json:"push_friends"`
	PushGroups       bool                 `json:"push_groups"`
	PushEvents       bool                 `json:"push_events"`
	PushBirthdays    bool                 `json:"push_birthdays"`
}

// --- CREATE REQ ---
type CreateUserNotificationSettingReq struct {
	UserID      string              `json:"user_id" binding:"required,uuid"`
	Settings    *GeneralSettingsReq `json:"settings" binding:"required"`
	Name        string              `json:"name" binding:"required"`
	Avatar      string              `bson:"avatar" json:"avatar"`
	DateOfBirth string              `json:"date_of_birth"`
	FCMTokens   []*FCMTokenReq      `json:"fcm_tokens"`
}

// --- UPDATE REQ ---
type UpdateUserNotificationSettingReq struct {
	DateOfBirth string              `json:"date_of_birth,omitempty"`
	Name        string              `json:"name,omitempty"`
	Avatar      string              `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Settings    *GeneralSettingsReq `json:"settings,omitempty"`
	FCMTokens   []*FCMTokenReq      `json:"fcm_tokens,omitempty"`
}
