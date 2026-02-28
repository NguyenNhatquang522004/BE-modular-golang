package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type ConversationAvatarReq struct {
	URL string `json:"url"`
}

type ConversationPermissionsReq struct {
	SendMessage *sharedEnums.PrivacyScope `json:"send_message"`
	AddMember   *sharedEnums.PrivacyScope `json:"add_member"`
}

type ConversationThemeReq struct {
	Color         string `json:"color"`
	Emoji         string `json:"emoji"`
	BackgroundURL string `json:"background_url"`
}

type LastMessageCacheReq struct {
	MessageID string                `json:"message_id"`
	Content   string                `json:"content"`
	SenderID  string                `json:"sender_id"`
	Type      sharedEnums.MediaType `json:"type"`
	CreatedAt time.Time             `json:"created_at"`
}

type ConversationReq struct {
	ID               string                        `json:"id"`
	Type             *enum.ConversationType        `json:"type" binding:"required"`
	Scope            *enum.ConversationScope       `json:"scope" binding:"required"`
	Status           *sharedEnums.ProcessingStatus `json:"status" binding:"required"`
	Name             string                        `json:"name,omitempty"`
	Avatar           *ConversationAvatarReq        `json:"avatar,omitempty"`
	CreatorID        string                        `json:"creator_id" binding:"required"`
	OwnerID          string                        `json:"owner_id" binding:"required"`
	RelatedGroupID   *string                       `json:"related_group_id,omitempty"` // Pointer string
	Permissions      ConversationPermissionsReq    `json:"permissions"`
	Theme            *ConversationThemeReq         `json:"theme,omitempty"`
	LastMessage      *LastMessageCacheReq          `json:"last_message,omitempty"`
	ParticipantCount int                           `json:"participant_count"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
}
