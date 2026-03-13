package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type ConversationAvatarRes struct {
	URL string `json:"url"`
}

type ConversationPermissionsRes struct {
	SendMessage sharedEnums.PrivacyScope `json:"send_message"`
	AddMember   sharedEnums.PrivacyScope `json:"add_member"`
}

type ConversationThemeRes struct {
	Color         string `json:"color"`
	Emoji         string `json:"emoji"`
	BackgroundURL string `json:"background_url"`
}

type LastMessageCacheRes struct {
	MessageID string                `json:"message_id"`
	Content   string                `json:"content"`
	SenderID  string                `json:"sender_id"`
	Type      sharedEnums.MediaType `json:"type"`
	CreatedAt time.Time             `json:"created_at"`
}

type ConversationRes struct {
	ID               string                        `json:"id"`
	Type             sharedEnums.ConversationType  `json:"type"`
	Scope            sharedEnums.ConversationScope `json:"scope"`
	Status           sharedEnums.ProcessingStatus  `json:"status"`
	Name             string                        `json:"name,omitempty"`
	Avatar           *ConversationAvatarRes        `json:"avatar,omitempty"`
	CreatorID        string                        `json:"creator_id"`
	OwnerID          string                        `json:"owner_id"`
	RelatedGroupID   *string                       `json:"related_group_id,omitempty"`
	RelatedChannelID *string                       `json:"related_channel_id,omitempty"`
	Permissions      ConversationPermissionsRes    `json:"permissions"`
	Theme            *ConversationThemeRes         `json:"theme,omitempty"`
	LastMessage      *LastMessageCacheRes          `json:"last_message,omitempty"`
	ParticipantCount int                           `json:"participant_count"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
}
