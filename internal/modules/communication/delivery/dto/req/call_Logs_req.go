package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
)

type CallLogReq struct {
	ID              string           `json:"id"` // Dùng string để Client truyền Hex string
	ConversationID  string           `json:"conversation_id" binding:"required"`
	CallerID        string           `json:"caller_id" binding:"required"`
	Participants    []string         `json:"participants" binding:"required"`
	Type            *enum.CallType   `json:"type" binding:"required"`
	Status          *enum.CallStatus `json:"status" binding:"required"`
	StartedAt       time.Time        `json:"started_at" binding:"required"`
	EndedAt         *time.Time       `json:"ended_at,omitempty"`
	DurationSeconds int              `json:"duration_seconds"`
	IsGroupCall     bool             `json:"is_group_call"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
}
