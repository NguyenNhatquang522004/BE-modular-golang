package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type ConversationParticipantRes struct {
	ID                string               `json:"id"`
	ConversationID    string               `json:"conversation_id"`
	UserID            string               `json:"user_id"`
	Role              sharedEnums.RoleType `json:"role"`
	Nickname          string               `json:"nickname"`
	LastSeenAt        time.Time            `json:"last_seen_at"`
	LastSeenMessageID string               `json:"last_seen_message_id"`
	MuteUntil         *time.Time           `json:"mute_until,omitempty"`
	IsArchived        bool                 `json:"is_archived"`
	ClearHistoryAt    *time.Time           `json:"clear_history_at,omitempty"`
	JoinedAt          time.Time            `json:"joined_at"`
	AddedByUserID     string               `json:"added_by_user_id"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	DeletedAt         *time.Time           `json:"deleted_at,omitempty"`
}