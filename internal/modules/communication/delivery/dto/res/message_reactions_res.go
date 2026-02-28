// File: res/message_reaction_res.go
package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// MessageReactionRes trả về thông tin đầy đủ của MessageReaction.
type MessageReactionRes struct {
	ConversationID string                   `json:"conversation_id"`
	MessageID      gocql.UUID               `json:"message_id"`
	UserID         gocql.UUID               `json:"user_id"`
	ReactionCode   sharedEnums.ReactionCode `json:"reaction_code"`
	CreatedAt      time.Time                `json:"created_at"`
}