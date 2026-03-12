package communicationEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

type MessagePayload struct {
	ConversationID   string                `json:"conversation_id"`
	Bucket           int                   `json:"bucket"`
	MessageID        gocql.UUID            `json:"message_id"`
	SenderID         gocql.UUID            `json:"sender_id" validate:"required"`
	Type             sharedEnums.MediaType `json:"type" validate:"required"`
	Content          string                `json:"content"`
	Attachments      []string              `json:"attachments"`
	IsEdited         bool                  `json:"is_edited"`
	ReplyToMessageID *gocql.UUID           `json:"reply_to_message_id,omitempty"`
	StoryRefID       *gocql.UUID           `json:"story_ref_id,omitempty"`
	IsRevoked        bool                  `json:"is_revoked"`
	CreatedAt        time.Time             `json:"created_at"`
	EventType        constants.EventType   `json:"event_type"`
}

type MessageStatsPayload struct {
	ConversationID string                   `json:"conversation_id"`
	Bucket         int                      `json:"bucket"`
	UserID         string                   `json:"user_id"`
	MessageID      string                   `json:"message_id"`
	ReactionCode   sharedEnums.ReactionCode `json:"reaction_code"`
	EventType      constants.EventType      `json:"event_type"`
}

type DeletePrivateConversationGroupPayload struct {
	TargetID string `json:"target_id"`
}
