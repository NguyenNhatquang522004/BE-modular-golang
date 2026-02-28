// File: req/message_reaction_req.go
package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// MessageReactionReq chứa 100% các trường từ Entity MessageReaction.
type MessageReactionReq struct {
	ConversationID string                   `json:"conversation_id" validate:"required"`
	MessageID      gocql.UUID               `json:"message_id" validate:"required"`
	UserID         gocql.UUID               `json:"user_id" validate:"required"`
	ReactionCode   sharedEnums.ReactionCode `json:"reaction_code" validate:"required"`
	CreatedAt      time.Time                `json:"created_at"`
}