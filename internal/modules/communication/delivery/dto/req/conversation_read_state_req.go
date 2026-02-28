// File: req/conversation_read_state_req.go
package req

import (
	"time"

	"github.com/gocql/gocql"
)

// ConversationReadStateReq chứa 100% các trường từ Entity ConversationReadState.
type ConversationReadStateReq struct {
	ConversationID    string     `json:"conversation_id" validate:"required"`
	UserID            gocql.UUID `json:"user_id" validate:"required"`
	LastReadMessageID gocql.UUID `json:"last_read_message_id" validate:"required"`
	LastReadAt        time.Time  `json:"last_read_at"`
}