package interactionEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

type EntityReactionPayload struct {
	TargetID     string                     `json:"target_id" validate:"required"`
	UserID       gocql.UUID                 `json:"user_id" validate:"required"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code" validate:"required"`
	CreatedAt    time.Time                  `json:"created_at"`
	Type         constants.EventType        `json:"topic" validate:"required"`
}
type CommentStatsPayload struct {
	TargetID     string              `json:"target_id" validate:"required"`
	UserID       string              `json:"user_id" validate:"required"`
	ReplyCount   int                 `json:"reply_count"`
	MentionCount int                 `json:"mention_count"`
	Total        int                 `json:"total"`
	Like         int                 `json:"like"`
	Love         int                 `json:"love"`
	Wow          int                 `json:"wow"`
	Sad          int                 `json:"sad"`
	Angry        int                 `json:"angry"`
	Type         constants.EventType `json:"type" validate:"required"` // CREATED, UPDATED, DELETED
}
type BookmarkPayload struct {
	UserID         string                      `json:"user_id" validate:"required"`
	TargetID       string                      `json:"target_id" validate:"required"`
	TargetType     sharedEnums.SavedTargetType `json:"target_type" validate:"required"`
	Snapshot       SavedItemSnapshotPayload    `json:"snapshot"`
	CollectionName string                      `json:"collection_name"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	ContextType    sharedEnums.ContextType     `json:"context_type"`
	Type           constants.EventType         `json:"type" validate:"required"` // CREATED, DELETED
}
type SavedItemSnapshotPayload struct {
	AuthorName     string `bson:"author_name" json:"author_name"`
	ContentPreview string `bson:"content_preview" json:"content_preview"` // Cắt ngắn 100 ký tự đầu
	ThumbnailURL   string `bson:"thumbnail_url" json:"thumbnail_url"`
}
type EditTargetPayload struct {
	ID               string                       `json:"id,omitempty"` // Có thể có hoặc không, tùy vào việc client có truyền lên hay không
	TargetCollection sharedEnums.TargetCollection `json:"target_collection" validate:"required"`
	TargetID         string                       `json:"target_id" validate:"required"`
	EditedAt         time.Time                    `json:"edited_at"`
	EditorID         string                       `json:"editor_id" validate:"required"`
	Diff             LogDiffPayload               `bson:"diff" json:"diff"`
	IPAddress        string                       `bson:"ip_address" json:"ip_address"`
	UserAgent        string                       `bson:"user_agent" json:"user_agent"`
	CreatedAt        time.Time                    `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time                    `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time                   `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
type LogDiffPayload struct {
	// Nội dung text
	OldContent string `bson:"old_content" json:"old_content"`
	NewContent string `bson:"new_content" json:"new_content"`

	// Media (Ảnh/Video)
	// Dùng Pointer để nếu không sửa ảnh thì field này null (omitempty) -> Tiết kiệm ổ cứng
	OldMedia *MediaSnapshotPayload `bson:"old_media,omitempty" json:"old_media,omitempty"`
	NewMedia *MediaSnapshotPayload `bson:"new_media,omitempty" json:"new_media,omitempty"`
}

// MediaSnapshot lưu trạng thái file media tại thời điểm sửa
type MediaSnapshotPayload struct {
	Type string `bson:"type" json:"type"` // "image", "video", "gif"
	URL  string `bson:"url" json:"url"`
}
type DeleteInteractionRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
