package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type CallLogRes struct {
	ID              string                 `json:"id"`
	ConversationID  string                 `json:"conversation_id"`
	CallerID        string                 `json:"caller_id"`
	Participants    []string               `json:"participants"`
	Type            sharedEnums.CallType   `json:"type"`
	Status          sharedEnums.CallStatus `json:"status"`
	StartedAt       time.Time              `json:"started_at"`
	EndedAt         *time.Time             `json:"ended_at,omitempty"`
	DurationSeconds int                    `json:"duration_seconds"`
	IsGroupCall     bool                   `json:"is_group_call"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
}
