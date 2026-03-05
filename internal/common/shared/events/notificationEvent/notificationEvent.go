package notificationEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type NotificationChangePayload struct {
	UserID      string                  `json:"user_id" binding:"required,uuid"`
	Settings    *GeneralSettingsPayload `json:"settings" binding:"required"`
	Name        string                  `json:"name" binding:"required"`
	Avatar      string                  `bson:"avatar" json:"avatar"`
	DateOfBirth string                  `json:"date_of_birth"`
	FCMTokens   []*FCMTokenPayload      `json:"fcm_tokens"`
}
type GeneralSettingsPayload struct {
	PushEnabled      bool                        `json:"push_enabled"`
	EmailFrequency   *sharedEnums.EmailFrequency `json:"email_frequency" binding:"required"`
	PushInteractions bool                        `json:"push_interactions"`
	PushFriends      bool                        `json:"push_friends"`
	PushGroups       bool                        `json:"push_groups"`
	PushEvents       bool                        `json:"push_events"`
	PushBirthdays    bool                        `json:"push_birthdays"`
}
type FCMTokenPayload struct {
	Token     string    `json:"token"`
	DeviceID  string    `json:"device_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationPayload struct {
	UserID           string                          `json:"user_id"`
	PostID           string                          `json:"post_id"`
	CommentID        string                          `json:"comment_id"`
	GroupID          string                          `json:"group_id"`
	CommentIDReply   string                          `json:"comment_id_reply"`
	SystemAlertID    string                          `json:"system_alert_id"`
	FriendAcceptID   string                          `json:"friend_accept_id"`
	BirthdayID       string                          `json:"birthday_user_id"`
	TypeNotification []*sharedEnums.NotificationType `json:"type_notification"`
	SendUser         []string                        `json:"send_user"`
	SendGroup        []string                        `json:"send_group"`
	SendPost         []string                        `json:"send_post"`
	Data             map[string]interface{}          `json:"data"`
	PreviewData      string                          `json:"preview_data"`
}
type NotifPostLikePayload struct {
}

type NotifFriendRequestPayload struct{}
type NotifFriendAcceptPayload struct{}
type NotifGroupInvitePayload struct{}
type NotifSystemAlertPayload struct{}
type NotifMentionPayload struct{}
type NotifFollowerPayload struct{}
type NotifFriendPayload struct{}
