package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type CreateCommentReq struct {
	TargetID        string           `json:"target_id" binding:"required"`
	UserID          string           `json:"user_id" binding:"required"`
	AssetID         *string          `json:"asset_id,omitempty"`
	Content         string           `json:"content"`
	Media           *CommentMediaReq `json:"media,omitempty"`
	Mentions        []string         `json:"mentions,omitempty"`
	ParentCommentID *string          `json:"parent_comment_id,omitempty"`
	RootCommentID   *string          `json:"root_comment_id,omitempty"`

	Status          *sharedEnums.ProcessingStatus `json:"status,omitempty"`
	HiddenMetadata  *HiddenMetadataReq            `json:"hidden_metadata,omitempty"`
	DeletedMetadata *DeletedMetadataReq           `json:"deleted_metadata,omitempty"`
	ReportCount     int                           `json:"report_count,omitempty"`
	Reactions       *CommentReactionsReq          `json:"reactions,omitempty"` // Đã tách riêng Req
	ReplyCount      int                           `json:"reply_count,omitempty"`
	MentionCount    int                           `json:"mention_count,omitempty"`
	IsEdited        bool                          `json:"is_edited,omitempty"`
	LastEditedAt    *time.Time                    `json:"last_edited_at,omitempty"`
	CreatedAt       *time.Time                    `json:"created_at,omitempty"`
	UpdatedAt       *time.Time                    `json:"updated_at,omitempty"`
}

type UpdateCommentReq struct {
	Content         *string                       `json:"content,omitempty"`
	Media           *CommentMediaReq              `json:"media,omitempty"`
	Mentions        *[]string                     `json:"mentions,omitempty"`
	Status          *sharedEnums.ProcessingStatus `json:"status,omitempty"`
	HiddenMetadata  *HiddenMetadataReq            `json:"hidden_metadata,omitempty"`
	DeletedMetadata *DeletedMetadataReq           `json:"deleted_metadata,omitempty"`
	ReportCount     *int                          `json:"report_count,omitempty"`
	Reactions       *CommentReactionsReq          `json:"reactions,omitempty"` // Đã tách riêng Req
	ReplyCount      *int                          `json:"reply_count,omitempty"`
	MentionCount    *int                          `json:"mention_count,omitempty"`
	IsEdited        *bool                         `json:"is_edited,omitempty"`
	LastEditedAt    *time.Time                    `json:"last_edited_at,omitempty"`
}

type CommentMediaReq struct {
	Type        *sharedEnums.MediaType `json:"type" binding:"required"`
	URL         string                 `json:"url" binding:"required,url"`
	DisplayMeta DisplayMetaReq         `json:"display_meta"`
}

type DisplayMetaReq struct {
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Duration  float64 `bson:"duration,omitempty" json:"duration,omitempty"`     // Dành cho video/audio
	SizeBytes int64   `bson:"size_bytes,omitempty" json:"size_bytes,omitempty"` // Dành cho tất cả loại media
	MimeType  string  `bson:"mime_type,omitempty" json:"mime_type,omitempty"`   // Dành cho tất cả loại media
}

type HiddenMetadataReq struct {
	IsHidden       bool      `json:"is_hidden"`
	HiddenAt       time.Time `json:"hidden_at"`
	HiddenByUserID string    `json:"hidden_by_user_id"`
	Reason         string    `json:"reason"`
	IsGhostBanned  bool      `json:"is_ghost_banned"`
}

type DeletedMetadataReq struct {
	DeletedAt       time.Time `json:"deleted_at"`
	DeletedByUserID string    `json:"deleted_by_user_id"`
}

type CommentReactionsReq struct {
	Total int `json:"total"`
	Like  int `json:"like"`
	Love  int `json:"love"`
	Haha  int `json:"haha"`
	Wow   int `json:"wow"`
	Sad   int `json:"sad"`
	Angry int `json:"angry"`
}
