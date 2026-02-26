package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
)

type CommentRes struct {
	ID              string              `json:"id"`
	PostID          string              `json:"post_id"`
	UserID          string              `json:"user_id"`
	AssetID         *string             `json:"asset_id,omitempty"`
	Content         string              `json:"content"`
	Media           *CommentMediaRes    `json:"media,omitempty"`
	Mentions        []string            `json:"mentions,omitempty"`
	ParentCommentID *string             `json:"parent_comment_id,omitempty"`
	RootCommentID   *string             `json:"root_comment_id,omitempty"`
	Status          *enum.CommentStatus `json:"status"`

	HiddenMetadata  *HiddenMetadataRes   `json:"hidden_metadata,omitempty"`
	DeletedMetadata *DeletedMetadataRes  `json:"deleted_metadata,omitempty"`
	ReportCount     int                  `json:"report_count"`
	Reactions       *CommentReactionsRes `json:"reactions,omitempty"` // Đã tách riêng Res
	ReplyCount      int                  `json:"reply_count"`
	MentionCount    int                  `json:"mention_count"`

	IsEdited     bool       `json:"is_edited"`
	LastEditedAt *time.Time `json:"last_edited_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CommentMediaRes struct {
	Type        enum.CommentMediaType `json:"type"`
	URL         string                `json:"url"`
	DisplayMeta DisplayMetaRes        `json:"display_meta"`
}

type DisplayMetaRes struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type HiddenMetadataRes struct {
	IsHidden       bool      `json:"is_hidden"`
	HiddenAt       time.Time `json:"hidden_at"`
	HiddenByUserID string    `json:"hidden_by_user_id"`
	Reason         string    `json:"reason"`
	IsGhostBanned  bool      `json:"is_ghost_banned"`
}

type DeletedMetadataRes struct {
	DeletedAt       time.Time `json:"deleted_at"`
	DeletedByUserID string    `json:"deleted_by_user_id"`
}

type CommentReactionsRes struct {
	Total int `json:"total"`
	Like  int `json:"like"`
	Love  int `json:"love"`
	Haha  int `json:"haha"`
	Wow   int `json:"wow"`
	Sad   int `json:"sad"`
	Angry int `json:"angry"`
}
