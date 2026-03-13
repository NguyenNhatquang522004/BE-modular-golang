package communicationEvent

import (
	"time"

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
	DeleteAll      bool                     `json:"delete_all,omitempty"` // Cờ để xóa tất cả các reaction của user này trên message này (nếu có)
}

type DeletePrivateConversationGroupPayload struct {
	TargetID string `json:"target_id"`
}

type CreateConversationPayload struct {
	ConversationID        string `json:"conversation_id" validate:"required"` // Client-side generated UUID
	UserCreatorAndOwnerID string `json:"user_id" validate:"required"`         // ID của người tạo cuộc hội thoại, dùng để add vào participant_ids bắt buộc
	UserCreatorName       string `json:"user_name" validate:"required"`       // Tên của người tạo cuộc hội thoại, dùng để làm nickname khi tạo participant record
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
	ParticipantIDs []participantInfo `json:"participant_ids" binding:"required,min=1"`
}
type participantInfo struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname,omitempty"`
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
	ConversationID string `json:"conversation_id" validate:"required"` // Client-side generated UUID
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

type CreateParticipantPayload struct {
	// ID cuộc hội thoại là bắt buộc
	ConversationID string `json:"conversation_id" validate:"required,mongodb"`

	// ID người dùng được thêm (Postgres UUID string)
	UserID string `json:"user_id" validate:"required,uuid4"`

	// Người thực hiện hành động thêm (Lấy từ Token/Context)
	AddedByUserID string `json:"added_by_user_id" validate:"required,uuid4"`

	// Vai trò mặc định thường là 'member', nhưng cho phép chỉ định nếu là Admin tạo nhóm
	Role sharedEnums.RoleType `json:"role" validate:"required,oneof=admin member"`

	// Biệt danh có thể có hoặc không
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
}
type DeleteParticipantPayload struct {
	ConversationID string `json:"conversation_id" validate:"required,mongodb"`
	UserID         string `json:"user_id" validate:"required,uuid4"`
	DeleteAll      bool   `json:"delete_all,omitempty"` // Cờ để xóa tất cả các bản ghi liên quan đến user này trong conversation (nếu có)
}
type UpdateParticipantPayload struct {
	ConversationID string `json:"conversation_id" validate:"required,mongodb"`
	UserID         string `json:"user_id" validate:"required,uuid4"`
	// Nickname mới (để trống nếu muốn xóa biệt danh)
	Nickname *string `json:"nickname" validate:"omitempty,max=50"`

	// Cập nhật vai trò (thường do Admin thực hiện)
	Role *sharedEnums.RoleType `json:"role" validate:"omitempty,oneof=admin member"`

	// Trạng thái lưu trữ cuộc hội thoại
	IsArchived *bool `json:"is_archived" validate:"omitempty"`

	// Tắt thông báo:
	// - Gửi thời gian cụ thể trong tương lai để mute.
	// - Gửi null (nil) để bật lại thông báo.
	MuteUntil *time.Time `json:"mute_until" validate:"omitempty"`
}
type CreateCallLogPayload struct {
	// Để string thay vì primitive.ObjectID trong Payload để parse JSON an toàn, sau đó map sang ObjectID ở Service
	ConversationID string `json:"conversation_id" validate:"required,mongodb"`

	// UUID của người gọi, có thể lấy từ JWT Token ở Controller, nhưng nếu bắt client gửi thì validate UUID
	CallerID string `json:"caller_id" validate:"required,uuid"`

	// Bắt buộc phải có ít nhất 2 người trong cuộc gọi, kiểm tra từng phần tử phải là chuẩn UUID
	Participants []string `json:"participants" validate:"required,min=2,dive,uuid"`

	// Type của cuộc gọi, cần validate đúng các giá trị enum cho phép
	Type sharedEnums.CallType `json:"type" validate:"required,oneof=voice video"`

	// Phân biệt call nhóm hay cá nhân
	IsGroupCall bool `json:"is_group_call"`
}
type UpdateCallLogPayload struct {
	ID string `json:"id" validate:"required,mongodb"` // ID của CallLog cần cập nhật, dùng để tìm bản ghi trong DB
	// Cập nhật trạng thái cuối cùng của cuộc gọi
	Status sharedEnums.CallStatus `json:"status" validate:"required,oneof=missed ended rejected"`

	// Thời điểm kết thúc do client báo lên (hoặc Server tự tính bằng time.Now())
	EndedAt time.Time `json:"ended_at" validate:"required"`

	// Thời lượng cuộc gọi tính bằng giây. Không thể là số âm.
	DurationSeconds int `json:"duration_seconds" validate:"gte=0"`
}

type DeleteCallLogPayload struct {
	ID             string `json:"id" validate:"required,mongodb"` // ID của CallLog cần xóa, dùng để tìm bản ghi trong DB
	ConversationID string `json:"conversation_id" validate:"required,mongodb"`
	DeleteAll      bool   `json:"delete_all,omitempty"` // Cờ để xóa tất cả các bản ghi liên quan đến conversation này (nếu có)
}

type ConversationStatsPayload struct {
	UserID            *string                   `json:"user_id"`
	ConversationID    string                   `json:"conversation_id"`
	ParticipantCount  *int                     `json:"participant_count"`
	LastMessage       *LastMessageCachePayload `json:"last_message"`
	LastSeenAt        *time.Time               `json:"last_seen_at"`
	LastSeenMessageID *string                  `json:"last_seen_message_id"`
	EventType         constants.EventType      `json:"event_type"`
}
type LastMessageCachePayload struct {
	// MessageID từ Cassandra (TimeUUID) -> Lưu String
	MessageID string `bson:"message_id" json:"message_id"`

	Content string `bson:"content" json:"content"`

	// SenderID từ Postgres (UUID) -> Lưu String
	SenderID string `bson:"sender_id" json:"sender_id"`

	Type sharedEnums.MediaType `bson:"type" json:"type"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
