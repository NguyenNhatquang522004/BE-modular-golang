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
type DeleteInteractionRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}
