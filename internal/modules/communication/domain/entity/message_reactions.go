package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

const (
	TableMessages = "messages"
)

// MessageReaction đại diện cho bảng 'message_reactions' trong Cassandra.
// Table này tối ưu cho việc đếm và liệt kê reaction của 1 tin nhắn.
type MessageReaction struct {
	// =========================================================================
	// 1. COMPOSITE PARTITION KEY
	// =========================================================================
	// Cả 2 trường này hợp lại xác định Node lưu trữ dữ liệu.

	// ID cuộc hội thoại (String - Mongo ID)
	ConversationID string `cql:"conversation_id" json:"conversation_id"`

	// ID tin nhắn (TimeUUID)
	MessageID gocql.UUID `cql:"message_id" json:"message_id"`

	// =========================================================================
	// 2. CLUSTERING KEY
	// =========================================================================

	// ID người thả reaction (Postgres UUID)
	// Đóng vai trò Clustering Key để đảm bảo 1 User chỉ thả 1 Reaction trên 1 tin nhắn.
	UserID gocql.UUID `cql:"user_id" json:"user_id"`

	// =========================================================================
	// 3. DATA FIELDS
	// =========================================================================

	// Enum Code: heart, haha, sad...
	ReactionCode sharedEnums.ReactionCode `cql:"reaction_code" json:"reaction_code"`

	CreatedAt time.Time `cql:"created_at" json:"created_at"`
}

func (Message) TableName() string {
	return TableMessages
}
