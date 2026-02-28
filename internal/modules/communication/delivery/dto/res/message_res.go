// File: res/message_res.go
package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// MessageRes trả về thông tin đầy đủ của Message.
type MessageRes struct {
	ConversationID   string                `json:"conversation_id"`
	Bucket           int                   `json:"bucket"`
	MessageID        gocql.UUID            `json:"message_id"`
	SenderID         gocql.UUID            `json:"sender_id"`
	Type             sharedEnums.MediaType `json:"type"`
	Content          string                `json:"content"`
	Attachments      []string              `json:"attachments"`
	IsEdited         bool                  `json:"is_edited"`
	ReplyToMessageID *gocql.UUID           `json:"reply_to_message_id,omitempty"`
	StoryRefID       *gocql.UUID           `json:"story_ref_id,omitempty"`
	IsRevoked        bool                  `json:"is_revoked"`
	CreatedAt        time.Time             `json:"created_at"`
}