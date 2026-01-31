package domain

import (
	"time"

	"github.com/gocql/gocql"
)

// ConversationReadState đại diện cho bảng 'conversation_read_state' trong Cassandra.
// Bảng này lưu trạng thái "Đã xem" (Read Receipt) của từng thành viên trong từng nhóm chat.
type ConversationReadState struct {
	// =========================================================================
	// 1. PRIMARY KEY (COMPOSITE)
	// =========================================================================

	// PARTITION KEY
	// Gom trạng thái đọc của 1 nhóm vào chung 1 node (thường thì partition key chỉ cần conversation_id là đủ).
	ConversationID string `cql:"conversation_id" json:"conversation_id"`

	// CLUSTERING KEY
	// Định danh user nào đang đọc.
	// Primary Key = (conversation_id, user_id) giúp query cực nhanh:
	// "User A đã đọc đến đâu trong nhóm X?"
	UserID gocql.UUID `cql:"user_id" json:"user_id"`

	// =========================================================================
	// 2. DATA FIELDS
	// =========================================================================

	// ID tin nhắn mới nhất mà user đã nhìn thấy (TimeUUID).
	// Dùng ID này so sánh với ID tin nhắn mới nhất của nhóm để tính số unread.
	LastReadMessageID gocql.UUID `cql:"last_read_message_id" json:"last_read_message_id"`

	// Thời điểm user thực hiện hành động đọc (Scroll xuống cuối).
	LastReadAt time.Time `cql:"last_read_at" json:"last_read_at"`
}

// TableName trả về tên bảng trong Cassandra
func (ConversationReadState) TableName() string {
	return "conversation_read_state"
}