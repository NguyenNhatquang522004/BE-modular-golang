// File: res/conversation_read_state_res.go
package res

import (
	"time"

	"github.com/gocql/gocql"
)

// ConversationReadStateRes trả về thông tin trạng thái đọc của người dùng.
type ConversationReadStateRes struct {
	ConversationID    string     `json:"conversation_id"`
	UserID            gocql.UUID `json:"user_id"`
	LastReadMessageID gocql.UUID `json:"last_read_message_id"`
	LastReadAt        time.Time  `json:"last_read_at"`
}