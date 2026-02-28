package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

const (
	TableMessageReactions = "message_reactions"
)

// Message đại diện cho bảng 'messages' trong Cassandra.
// Thiết kế theo mô hình "Wide Rows" với Time-Bucketing pattern.
type Message struct {
	// =========================================================================
	// 1. PRIMARY KEY COMPONENTS
	// =========================================================================

	// PARTITION KEY 1: ID cuộc hội thoại (String vì gốc là Mongo ObjectID)
	ConversationID string `cql:"conversation_id" json:"conversation_id"`

	// PARTITION KEY 2: Bucket thời gian (VD: 202310 - Tháng 10/2023)
	// Giúp chia nhỏ partition, tránh row quá lớn (>100MB) gây chậm Cassandra.
	Bucket int `cql:"bucket" json:"bucket"`

	// CLUSTERING KEY: TimeUUID
	// Vừa là ID duy nhất, vừa chứa timestamp để sắp xếp tin nhắn.
	MessageID gocql.UUID `cql:"message_id" json:"message_id"`

	// =========================================================================
	// 2. MESSAGE DATA
	// =========================================================================

	// SenderID: User ID từ Postgres (UUID)
	SenderID gocql.UUID `cql:"sender_id" json:"sender_id"`

	// Enum Type: text, image, video...
	Type sharedEnums.MediaType `cql:"type" json:"type"`

	Content     string   `cql:"content" json:"content"`
	Attachments []string `cql:"attachments" json:"attachments"` // LIST<TEXT>

	// =========================================================================
	// 3. REFERENCES (Nullable -> Pointer)
	// =========================================================================
	IsEdited bool `cql:"is_edited" json:"is_edited"` // Tin nhắn đã bị chỉnh sửa hay chưa?

	// Reply tin nhắn nào trong cùng conversation?
	ReplyToMessageID *gocql.UUID `cql:"reply_to_message_id" json:"reply_to_message_id"`

	// Reply story nào?
	// Lưu ý: Nếu StoryID là Mongo ObjectId, bạn cần đổi SQL thành TEXT và Go thành *string.
	// Ở đây tôi tuân thủ SQL của bạn là UUID.
	StoryRefID *gocql.UUID `cql:"story_ref_id" json:"story_ref_id"`

	// =========================================================================
	// 4. METADATA
	// =========================================================================

	IsRevoked bool      `cql:"is_revoked" json:"is_revoked"` // Thu hồi tin nhắn
	CreatedAt time.Time `cql:"created_at" json:"created_at"`
}

// TableName trả về tên bảng trong Cassandra
func (MessageReaction) TableName() string {
	return TableMessageReactions
}
