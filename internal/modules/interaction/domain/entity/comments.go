package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
	collectionnamComment = "Comment"
)
type Comment struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Reference bài viết gốc (Mongo ID)
	// Index: { post_id: 1, created_at: 1 } -> Load comment của bài viết
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`

	// Người comment (Postgres UUID) -> String
	UserID string `bson:"user_id" json:"user_id"`

	// 1. NGỮ CẢNH (Context)
	// Nếu comment vào 1 ảnh cụ thể trong bài Multi-Photo
	AssetID *primitive.ObjectID `bson:"asset_id,omitempty" json:"asset_id,omitempty"`

	// 2. NỘI DUNG
	Content string        `bson:"content" json:"content"`
	Media   *CommentMedia `bson:"media,omitempty" json:"media,omitempty"`

	// Danh sách UserID (UUID String) được tag
	Mentions []string `bson:"mentions,omitempty" json:"mentions,omitempty"`

	// 3. CẤU TRÚC PHẢN HỒI (Nested Comments)
	// Dùng Pointer để có thể null (Level 1 thì parent/root = nil)
	ParentCommentID *primitive.ObjectID `bson:"parent_comment_id,omitempty" json:"parent_comment_id,omitempty"`
	RootCommentID   *primitive.ObjectID `bson:"root_comment_id,omitempty" json:"root_comment_id,omitempty"`

	// 4. TRẠNG THÁI & MODERATION
	Status enum.CommentStatus `bson:"status" json:"status"`

	// Pointer struct: Nếu không bị ẩn/xóa thì field này null -> Tiết kiệm data
	HiddenMetadata  *HiddenMetadata  `bson:"hidden_metadata,omitempty" json:"hidden_metadata,omitempty"`
	DeletedMetadata *DeletedMetadata `bson:"deleted_metadata,omitempty" json:"deleted_metadata,omitempty"`

	ReportCount int `bson:"report_count" json:"report_count"`

	// 5. METRICS
	Reactions    CommentReactions `bson:"reactions" json:"reactions"`
	ReplyCount   int              `bson:"reply_count" json:"reply_count"`
	MentionCount int              `bson:"mention_count" json:"mention_count"`

	// 6. LỊCH SỬ CHỈNH SỬA
	IsEdited     bool       `bson:"is_edited" json:"is_edited"`
	LastEditedAt *time.Time `bson:"last_edited_at,omitempty" json:"last_edited_at,omitempty"`

	// 7. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- MEDIA ---
type CommentMedia struct {
	Type        enum.CommentMediaType `bson:"type" json:"type"`
	URL         string                `bson:"url" json:"url"`
	DisplayMeta DisplayMeta           `bson:"display_meta" json:"display_meta"`
}

type DisplayMeta struct {
	Width  int `bson:"width" json:"width"`
	Height int `bson:"height" json:"height"`
}

// --- METADATA (Hidden/Deleted) ---
type HiddenMetadata struct {
	IsHidden bool      `bson:"is_hidden" json:"is_hidden"`
	HiddenAt time.Time `bson:"hidden_at" json:"hidden_at"`

	// UserID (UUID Postgres) thực hiện ẩn -> String
	HiddenByUserID string `bson:"hidden_by_user_id" json:"hidden_by_user_id"`

	Reason        string `bson:"reason" json:"reason"`
	IsGhostBanned bool   `bson:"is_ghost_banned" json:"is_ghost_banned"`
}

type DeletedMetadata struct {
	DeletedAt time.Time `bson:"deleted_at" json:"deleted_at"`

	// UserID (UUID Postgres) thực hiện xóa -> String
	DeletedByUserID string `bson:"deleted_by_user_id" json:"deleted_by_user_id"`
}

// --- REACTIONS ---
// Denormalization để hiển thị UI nhanh
type CommentReactions struct {
	Total int `bson:"total" json:"total"`
	Like  int `bson:"like" json:"like"`
	Love  int `bson:"love" json:"love"`
	Haha  int `bson:"haha" json:"haha"`
	Wow   int `bson:"wow" json:"wow"`
	Sad   int `bson:"sad" json:"sad"`
	Angry int `bson:"angry" json:"angry"`
}

func (Comment) CollectionnamComment() string {
	return collectionnamComment
}
