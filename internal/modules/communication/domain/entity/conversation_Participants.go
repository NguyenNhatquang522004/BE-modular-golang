package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionConversationParticipants = "ConversationParticipants"
)

// ConversationParticipant đại diện cho 1 thành viên trong 1 cuộc hội thoại.
// Đây là collection TRUNG GIAN quan trọng nhất để render danh sách Inbox.
type ConversationParticipant struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// --- 1. LINKING (KHÓA NGOẠI) ---

	// ID cuộc hội thoại (Mongo)
	// Index: { conversation_id: 1 } -> Lấy danh sách thành viên của 1 nhóm
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`

	// ID thành viên (Postgres UUID) -> String
	// Index: { user_id: 1, last_seen_at: -1 } -> Lấy danh sách Inbox của user (Sắp xếp theo hoạt động)
	UserID string `bson:"user_id" json:"user_id"`

	// --- 2. VAI TRÒ & ĐỊNH DANH ---
	Role     sharedEnums.RoleType `bson:"role" json:"role"`         // 'admin', 'member'
	Nickname string               `bson:"nickname" json:"nickname"` // Biệt danh trong nhóm này

	// --- 3. TRẠNG THÁI ĐỌC (READ STATUS) ---
	// Thời điểm cuối user mở box chat
	LastSeenAt time.Time `bson:"last_seen_at" json:"last_seen_at"`

	// ID tin nhắn cuối cùng đã đọc (Cassandra TimeUUID -> String)
	// Dùng để tính Unread Count: Count(Messages where ID > LastSeenMessageID)
	LastSeenMessageID string `bson:"last_seen_message_id" json:"last_seen_message_id"`

	// --- 4. CÀI ĐẶT CÁ NHÂN (PERSONAL SETTINGS) ---
	MuteUntil  *time.Time `bson:"mute_until,omitempty" json:"mute_until,omitempty"` //Cho phép tắt thông báo trong 1 giờ, 8 giờ hoặc mãi mãi. sửa cái này
	IsArchived bool       `bson:"is_archived" json:"is_archived"`                   // Lưu trữ (ẩn khỏi Inbox)

	// --- 5. TÍNH NĂNG XÓA LỊCH SỬ ---
	// Mốc thời gian xóa tin nhắn phía client.
	// Pointer để check null (Mặc định là nil - xem full lịch sử)
	ClearHistoryAt *time.Time `bson:"clear_history_at,omitempty" json:"clear_history_at,omitempty"`

	// --- 6. METADATA ---
	JoinedAt time.Time `bson:"joined_at" json:"joined_at"`

	// Người add user này vào nhóm (Postgres UUID -> String)
	AddedByUserID string `bson:"added_by_user_id" json:"added_by_user_id"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (ConversationParticipant) CollectionName() string {
	return CollectionConversationParticipants
}
