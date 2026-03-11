package mediaEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type ReplyStoryPayload struct {
	StoryID string `json:"story_id"`
	ReplyID int    `json:"reply_id"`
}

type CoutnerLiveStreamPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	Comments      int                 `json:"comments"`
	Views         int                 `json:"views"`
	EventType     constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}

type StartStopVideoLiveStreamPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	SegmentLen    int                 `json:"segment_len"`
	OwnerID       string              `json:"owner_id"`
	Name          string              `json:"name"`
	EventType     constants.EventType `json:"event_type"`
}

type ReactLiveStreamPayload struct {
	LiveSessionID string                     `json:"live_session_id"`
	UserID        string                     `json:"user_id"`
	Total         int                        `json:"total"`
	TargetType    sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode  sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt     time.Time                  `json:"created_at"`
	EventType     constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}
type ReactCounterReelPayload struct {
	ReelID    string              `json:"reel_id"`
	UserID    string              `json:"user_id"`
	Comments  int                 `json:"comments"`
	Saves     int                 `json:"saves"`
	Shares    int                 `json:"shares"`
	EventType constants.EventType `json:"event_type"` // "view" hoặc "unview"
}

type ReactReelPayload struct {
	ReelID       string                     `json:"reel_id"`
	UserID       string                     `json:"user_id"`
	Total        int                        `json:"total"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt    time.Time                  `json:"created_at"`
	EventType    constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}

type DeleteMediaRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}

// /
type CreateMediaAssetsPayload struct {
	// Chuyển toàn bộ primitive.ObjectID của MongoDB thành string
	// // Nếu client gửi lên có nghĩa là update, nếu không có nghĩa là create mới
	AssetsID string             `json:"assets_id,omitempty"`
	UserID   string             `json:"user_id"` // ID người upload (UUID từ Postgres)
	Items    []MediaItemPayload `json:"items"`
}

// --- 3. SUB-STRUCT: MEDIA ITEM ---
type MediaItemPayload struct {
	MediaID      string                `json:"media_id,omitempty"` // Nếu client gửi lên có nghĩa là update, nếu không có nghĩa là create mới
	PostID       string                `json:"post_id"`
	AlbumID      string                `json:"album_id,omitempty"`
	GroupID      string                `json:"group_id,omitempty"`
	CommentID    string                `json:"comment_id,omitempty"`
	PageID       string                `json:"page_id,omitempty"`
	MediaType    sharedEnums.MediaType `json:"media_type"` // Sử dụng Enum đã định nghĩa
	URL          string                `json:"url"`
	ThumbnailURL string                `json:"thumbnail_url"`
	Metadata     MetadataPayload       `json:"metadata"`
	Order        int                   `json:"order"`
	Hashtags     []string              `json:"hashtags,omitempty"`
	// Lưu ý: Nếu module Media KHÔNG quan tâm đến TaggedUser (chỉ quan tâm xử lý file/ảnh),
	// bạn hoàn toàn có thể lược bỏ TaggedUsers ở đây để giữ payload nhẹ (Thin Payload).
	// Dưới đây vẫn giữ lại để đảm bảo đủ 100% data như entity của bạn.
	TaggedUsers []TaggedUserPayload `json:"tagged_users,omitempty"`
}

// --- 4. SUB-STRUCT: METADATA ---
type MetadataPayload struct {
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes"`
	MimeType  string  `json:"mime_type"`
}

// --- 5. SUB-STRUCT: TAGGED USER ---
type TaggedUserPayload struct {
	UserID string  `json:"user_id"` // Đã là string từ Postgres UUID
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type UpdateMediaAssetsPayload struct {
	MediaID string `json:"media_id"`
	// Các trường có thể cập nhật
	URL          *string              `json:"url,omitempty"`
	ThumbnailURL *string              `json:"thumbnail_url,omitempty"`
	Order        *int                 `json:"order,omitempty"`
	Hashtags     *[]string            `json:"hashtags,omitempty"`
	Metadata     *MetadataPayload     `json:"metadata,omitempty"`
	TaggedUsers  *[]TaggedUserPayload `json:"tagged_users,omitempty"`
	Privacy      *MediaPrivacyPayload `json:"privacy,omitempty"`
}

type DeleteMediaAssetsPayload struct {
	MediaID string `json:"media_id"`
}
type MediaPrivacyPayload struct {
	// Level string hoặc dùng Enum PrivacyScope tái sử dụng
	Level            sharedEnums.PrivacyScope `bson:"level" json:"level"`
	InheritFromAlbum bool                     `bson:"inherit_from_album" json:"inherit_from_album"`
}
type DeleteMediaByTargetPayload struct {
	TargetID string `json:"target_id"`
}

type AlbumStatsPayload struct {
	UserID     string              `json:"user_id"`
	AlbumID    string              `json:"album_id"`
	AssetCount int                 `json:"asset_count"`
	Like       int                 `json:"like"`
	Love       int                 `json:"love"`
	Haha       int                 `json:"haha"`
	Wow        int                 `json:"wow"`
	Sad        int                 `json:"sad"`
	Angry      int                 `json:"angry"`
	EventType  constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}

type StoryStatsPayload struct {
	UserID          string                           `json:"user_id"`
	StoryID         string                           `json:"story_id"`
	Views           int                              `json:"views"`
	Like            int                              `json:"like"`
	Love            int                              `json:"love"`
	Haha            int                              `json:"haha"`
	Wow             int                              `json:"wow"`
	Sad             int                              `json:"sad"`
	Angry           int                              `json:"angry"`
	ReplyCount      int                              `json:"reply_count"`
	ViewsCount      int                              `json:"views_count"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`
	PollOptionIndex *int                             `json:"poll_option_index,omitempty"` // Chỉ có khi InteractionType là
	Content         string                           `json:"content,omitempty"`           // Chỉ có khi InteractionType là reaction hoặc comment
	EventType       constants.EventType              `json:"event_type"`                  // "increment" hoặc "decrement"
}

type ReelStatsPayload struct {
	UserID    string              `json:"user_id"`
	ReelID    string              `json:"reel_id"`
	Views     int                 `json:"views"`
	Like      int                 `json:"like"`
	Love      int                 `json:"love"`
	Haha      int                 `json:"haha"`
	Wow       int                 `json:"wow"`
	Sad       int                 `json:"sad"`
	Angry     int                 `json:"angry"`
	Shares    int                 `json:"shares"`
	Saves     int                 `json:"saves"`
	Comments  int                 `json:"comments"`
	EventType constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}

type LiveSessionStatsPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	UserID        string              `json:"user_id"`
	PeakViewers   int                 `json:"peak_viewers"`
	TotalViews    int                 `json:"total_views"`
	TotalComments int                 `json:"total_comments"`
	Like          int                 `json:"like"`
	Love          int                 `json:"love"`
	Haha          int                 `json:"haha"`
	Wow           int                 `json:"wow"`
	Sad           int                 `json:"sad"`
	Angry         int                 `json:"angry"`
	EventType     constants.EventType `json:"event_type"` // "view" hoặc "unview"
}

type LiveCommentPayload struct {
	UserID        string                   `json:"user_id"`
	StreamID      string                   `json:"stream_id"`
	CreatedAt     time.Time                `json:"created_at"`
	CommentID     string                   `json:"comment_id"`
	UserBadges    []*sharedEnums.UserBadge `json:"user_badges"`
	Content       string                   `json:"content"`
	IsPinned      bool                     `json:"is_pinned"`
	UserNickname  string                   `json:"user_nickname"`
	UserAvatarURL string                   `json:"user_avatar_url"`
	EventType     constants.EventType      `json:"event_type"`
}

// CreateStoryPayload đại diện cho payload Client gửi lên khi tạo Story mới
type CreateStoryPayload struct {
	Media    StoryMediaPayload     `json:"media" binding:"required"`
	Overlays []StoryOverlayPayload `json:"overlays,omitempty" binding:"dive"` // dive: validate từng phần tử trong mảng
	Privacy  StoryPrivacyPayload   `json:"privacy" binding:"required"`
	Settings StorySettingsPayload  `json:"settings" binding:"required"`
}

type StoryMediaPayload struct {
	URL          string                `json:"url" binding:"required,url"` // Bắt buộc phải là định dạng URL
	Type         sharedEnums.MediaType `json:"type" binding:"required"`
	Duration     float64               `json:"duration" binding:"gte=0"` // Lớn hơn hoặc bằng 0
	ThumbnailURL string                `json:"thumbnail_url" binding:"omitempty,url"`
	SizeBytes    int64                 `json:"size_bytes" binding:"gte=0"`
}

type StoryPrivacyPayload struct {
	Type sharedEnums.PrivacyScope `json:"type" binding:"required"`
	// Nếu dùng UUID cho Postgres, validate uuid ở đây
	AllowList []string `json:"allow_list,omitempty" binding:"omitempty,dive,uuid"`
	BlockList []string `json:"block_list,omitempty" binding:"omitempty,dive,uuid"`
}

type StorySettingsPayload struct {
	// Không dùng pointer ở Create vì ta cần force client gửi các config này
	AllowReply bool `json:"allow_reply"`
	AllowShare bool `json:"allow_share"`
}

type StoryOverlayPayload struct {
	Type     sharedEnums.OverlayType       `json:"type" binding:"required"`
	Position OverlayPositionPayload `json:"position" binding:"required"`
	Data     map[string]interface{} `json:"data" binding:"required"`
}

type OverlayPositionPayload struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
	Scale    float64 `json:"scale" binding:"gt=0"` // Scale phải lớn hơn 0
}
type UpdateStoryPayload struct {
	Privacy  *UpdateStoryPrivacyPayload  `json:"privacy,omitempty"`
	Settings *UpdateStorySettingsPayload `json:"settings,omitempty"`
}

type UpdateStoryPrivacyPayload struct {
	Type      *sharedEnums.PrivacyScope `json:"type,omitempty"`
	AllowList []string                  `json:"allow_list,omitempty" binding:"omitempty,dive,uuid"`
	BlockList []string                  `json:"block_list,omitempty" binding:"omitempty,dive,uuid"`
}

type UpdateStorySettingsPayload struct {
	// Bắt buộc phải dùng con trỏ (pointer) cho boolean trong Update DTO.
	// Nếu dùng bool thường, khi client không truyền `allow_reply`, Go sẽ tự hiểu là `false`.
	// Dùng `*bool` giúp ta check: nếu nó là nil -> client không muốn update trường này.
	AllowReply *bool `json:"allow_reply,omitempty"`
	AllowShare *bool `json:"allow_share,omitempty"`
}
