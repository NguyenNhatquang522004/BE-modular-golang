package communicationEvent

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

type AttachmentPayload struct {
	AssetID      string                `json:"asset_id"` // ID tạm hoặc ID đã tạo từ phía Upload Service
	URL          string                `json:"url" validate:"required"`
	ThumbnailURL string                `json:"thumbnail_url"`
	Type         sharedEnums.MediaType `json:"type" validate:"required"`
	Width        int                   `json:"width"`
	Height       int                   `json:"height"`
	SizeBytes    int64                 `json:"size_bytes"`
	MimeType     string                `json:"mime_type"`
}
type CreateMessagePayload struct {
	ConversationID string     `json:"conversation_id" validate:"required"`
	MessageID      string     `json:"message_id" validate:"required"` // Client-side generated UUID
	SenderID       gocql.UUID `json:"sender_id" validate:"required"`

	Type    sharedEnums.MediaType `json:"type" validate:"required"`
	Content string                `json:"content"`

	// Cải tiến: Danh sách các đối tượng Media chi tiết
	Attachments []AttachmentPayload `json:"attachments"`

	ReplyToMessageID *gocql.UUID `json:"reply_to_message_id,omitempty"`
	StoryRefID       *gocql.UUID `json:"story_ref_id,omitempty"`
}
type DeleteMessagePayload struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Bucket         int    `json:"bucket" validate:"required"`
	MessageID      string `json:"message_id" validate:"required"`
	DeleteAlll     bool   `json:"delete_all,omitempty"` // Cờ để xóa tất cả các phiên bản của message (nếu có)
}
type UpdateMessagePayload struct {
	// --- IDENTIFIERS ---
	ConversationID string `json:"conversation_id" validate:"required"`
	Bucket         int    `json:"bucket" validate:"required"`
	MessageID      string `json:"message_id" validate:"required"`

	// --- UPDATABLE FIELDS ---
	Content *string `json:"content,omitempty"`

	// Danh sách Attachment mới (bao gồm cả những cái cũ giữ lại và cái mới thêm vào)
	Attachments *[]AttachmentPayload `json:"attachments,omitempty"`

	// CẢI TIẾN: Danh sách các ID cần xóa bỏ hẳn khỏi MediaAsset hoặc đánh dấu Deleted
	RemoveAttachmentIDs []string `json:"remove_attachment_ids,omitempty"`

	IsRevoked *bool `json:"is_revoked,omitempty"`
}
type MessageStatsPayload struct {
	ConversationID string                   `json:"conversation_id"`
	Bucket         int                      `json:"bucket"`
	UserID         string                   `json:"user_id"`
	MessageID      string                   `json:"message_id"`
	ReactionCode   sharedEnums.ReactionCode `json:"reaction_code"`
	EventType      constants.EventType      `json:"event_type"`
}

type DeletePrivateConversationGroupPayload struct {
	TargetID string `json:"target_id"`
}

type CreateConversationPayload struct {
	ConversationID        string `json:"conversation_id" validate:"required"` // Client-side generated UUID
	UserCreatorAndOwnerID string `json:"user_id" validate:"required"`         // ID của người tạo cuộc hội thoại, dùng để add vào participant_ids bắt buộc
	// --- 1. CLASSIFICATION ---
	Type  sharedEnums.ConversationType  `json:"type" binding:"required"`  // Bắt buộc: 'private' hoặc 'group'
	Scope sharedEnums.ConversationScope `json:"scope" binding:"required"` // Bắt buộc: 'messenger' hoặc 'community_channel'

	// --- 2. GROUP INFO ---
	// Name không bắt buộc ở payload vì nếu là private (1-1) thì không cần tên.
	// Sẽ validate logic ở tầng Service: nếu Type == 'group' thì Name mới bắt buộc.
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar_url,omitempty"` // Nhận URL string cho gọn, Service sẽ map vào struct ConversationAvatar

	// --- 3. LINKING ---
	// Best Practice: Nhận String từ JSON, sau đó parse sang primitive.ObjectID ở tầng Service
	// Giải thích:
	// - RelatedGroupID: ID của Group tổng.
	// - RelatedChannelID: ID của Channel nằm trong Group tổng, chứa conversation này.
	RelatedGroupID   string `json:"related_group_id,omitempty"`
	RelatedChannelID string `json:"related_channel_id,omitempty"`

	// --- 4. SETTINGS ---
	Permissions *ConversationPermissionsPayload `json:"permissions,omitempty"`
	Theme       *ConversationThemePayload       `json:"theme,omitempty"`

	// --- LƯU Ý QUAN TRỌNG (BEST PRACTICE) ---
	// Dù Entity Conversation không trực tiếp lưu mảng UserID (do bạn dùng Collection/Table Members riêng),
	// nhưng khi TẠO cuộc hội thoại, client phải gửi lên danh sách những người được thêm vào.
	ParticipantIDs []string `json:"participant_ids" binding:"required,min=1"`
}

// ConversationPermissionsPayload payload cho permissions
type ConversationPermissionsPayload struct {
	SendMessage sharedEnums.PrivacyScope `json:"send_message" binding:"required"`
	AddMember   sharedEnums.PrivacyScope `json:"add_member" binding:"required"`
}

// ConversationThemePayload payload cho theme
type ConversationThemePayload struct {
	Color         string `json:"color,omitempty"`          // Ví dụ: #FF0000
	Emoji         string `json:"emoji,omitempty"`          // Ví dụ: 👍
	BackgroundURL string `json:"background_url,omitempty"` // URL ảnh nền
}
type UpdateConversationReq struct {
	// --- 1. GROUP INFO ---
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar_url,omitempty"` // Tương tự, dùng string pointer chứa URL

	// --- 2. SETTINGS ---
	Permissions *ConversationPermissionsPayload `json:"permissions,omitempty"`
	Theme       *ConversationThemePayload       `json:"theme,omitempty"`

	// --- 3. CLASSIFICATION / OWNERSHIP ---
	// Thông thường Type và Scope không cho phép sửa sau khi tạo.
	// Status có thể cập nhật (ví dụ: client muốn Archive/Xoá tạm thời đoạn chat)
	Status *sharedEnums.ProcessingStatus `json:"status,omitempty"`

	// Bổ sung tính năng chuyển nhượng quyền Owner (chỉ Owner hiện tại mới có quyền gọi payload có field này)
	OwnerID *string `json:"owner_id,omitempty"`
}
