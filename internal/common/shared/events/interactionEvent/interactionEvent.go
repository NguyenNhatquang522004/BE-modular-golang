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
	Type           constants.EventType         `json:"type" validate:"required"` // CREATED, DELETED
}
type SavedItemSnapshotPayload struct {
	AuthorName     string `bson:"author_name" json:"author_name"`
	ContentPreview string `bson:"content_preview" json:"content_preview"` // Cắt ngắn 100 ký tự đầu
	ThumbnailURL   string `bson:"thumbnail_url" json:"thumbnail_url"`
}

type DeleteInteractionRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
