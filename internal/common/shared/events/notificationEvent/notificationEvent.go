package notificationEvent

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

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
