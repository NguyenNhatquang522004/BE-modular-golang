package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type BlockRequest struct {
	ID            string                   `json:"id"`
	BlockerUserID string                   `json:"blocker_user_id"`
	BlockedUserID string                   `json:"blocked_user_id"`
	TypeBlock     []sharedEnums.Type_Block `json:"type_block"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	EventType     constants.EventType      `json:"event_type"`
}
