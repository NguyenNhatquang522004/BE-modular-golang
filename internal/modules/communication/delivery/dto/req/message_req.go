// File: req/message_req.go
package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// MessageReq chứa đầy đủ 100% các trường từ Entity Message.
type MessageReq struct {
	ConversationID   string                `json:"conversation_id" validate:"required"`
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
}