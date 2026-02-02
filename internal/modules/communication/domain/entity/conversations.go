package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
	CollectionConversations = "Conversations"
)
type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. CLASSIFICATION
	Type   enum.ConversationType   `bson:"type" json:"type"`     // 'private', 'group'
	Scope  enum.ConversationScope  `bson:"scope" json:"scope"`   // 'messenger', 'community_channel'
	Status enum.ConversationStatus `bson:"status" json:"status"` // 'active', 'pending'

	// 2. GROUP INFO
	// Tên nhóm (Null nếu là Private chat 1-1)
	Name   string              `bson:"name,omitempty" json:"name,omitempty"`
	Avatar *ConversationAvatar `bson:"avatar,omitempty" json:"avatar,omitempty"`

	// 3. OWNERSHIP & LINKING
	// UserID từ Postgres (UUID) -> String
	CreatorID string `bson:"creator_id" json:"creator_id"`
	OwnerID   string `bson:"owner_id" json:"owner_id"`

	// Link tới Module Groups (Nếu scope là community_channel)
	// Index: { related_group_id: 1 } -> Lấy danh sách kênh của 1 nhóm
	RelatedGroupID *primitive.ObjectID `bson:"related_group_id,omitempty" json:"related_group_id,omitempty"`

	// 4. SETTINGS
	// Permissions luôn có giá trị mặc định, không nên để pointer
	Permissions ConversationPermissions `bson:"permissions" json:"permissions"`

	// Theme có thể null (dùng mặc định của App)
	Theme *ConversationTheme `bson:"theme,omitempty" json:"theme,omitempty"`

	// 5. CACHING (Quan trọng cho Inbox List)
	// Pointer để check null (Nhóm mới tạo chưa chat gì)
	LastMessage *LastMessageCache `bson:"last_message,omitempty" json:"last_message,omitempty"`

	// Counter thành viên (Denormalization để không phải count bảng Members)
	ParticipantCount int `bson:"participant_count" json:"participant_count"`

	// 6. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// --- AVATAR ---
type ConversationAvatar struct {
	URL string `bson:"url" json:"url"`
}

// --- PERMISSIONS ---
type ConversationPermissions struct {
	SendMessage enum.PermissionLevel `bson:"send_message" json:"send_message"` // 'everyone', 'admin_only'
	AddMember   enum.PermissionLevel `bson:"add_member" json:"add_member"`
}

// --- THEME ---
type ConversationTheme struct {
	Color         string `bson:"color" json:"color"` // Hex code: #FF0000
	Emoji         string `bson:"emoji" json:"emoji"` // "👍"
	BackgroundURL string `bson:"background_url" json:"background_url"`
}

// --- LAST MESSAGE CACHE ---
// Denormalization: Giúp load Inbox nhanh mà không cần query Cassandra
type LastMessageCache struct {
	// MessageID từ Cassandra (TimeUUID) -> Lưu String
	MessageID string `bson:"message_id" json:"message_id"`

	Content string `bson:"content" json:"content"`

	// SenderID từ Postgres (UUID) -> Lưu String
	SenderID string `bson:"sender_id" json:"sender_id"`

	Type enum.MessageType `bson:"type" json:"type"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
func (Conversation) CollectionName() string {
	return CollectionConversations
}