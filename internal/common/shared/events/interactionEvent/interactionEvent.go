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
	Topic        constants.EventType        `json:"topic" validate:"required"`
}

type CounterPostPayload struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}

type CommentCountPayload struct {
	CommentID    string              `json:"comment_id,omitempty"` // Nếu có comment_id thì đếm reply của comment đó, không có thì đếm comment của post
	ReplyCount   int                 `json:"reply_count,omitempty"`
	MentionCount int                 `json:"mention_count,omitempty"`
	ReportCount  int                 `json:"report_count,omitempty"`
	EventType    constants.TopicName `json:"event_type,omitempty"`
}

type DeleteInteractionRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
