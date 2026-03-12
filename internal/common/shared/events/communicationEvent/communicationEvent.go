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
