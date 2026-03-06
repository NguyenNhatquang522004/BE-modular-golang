package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
)

type FollowerRequest struct {
	ID              string              `json:"id"`
	Follower_UserID string              `json:"follower_user_id"`
	Followed_UserID string              `json:"followed_user_id"`
	IsMuted         bool                `json:"is_muted"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	EventType       constants.EventType `json:"event_type"`
}
