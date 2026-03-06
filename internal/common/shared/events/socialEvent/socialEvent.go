package socialEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type DeleteSocialRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}

type FollowerUserPayload struct {
	ID             string              `json:"id"`
	FollowerUserID string              `json:"follower_user_id"`
	FollowedUserID string              `json:"followed_user_id"`
	IsMuted        bool                `json:"is_muted"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	EventType      constants.EventType `json:"event_type"`
}

type FriendUserPayload struct {
	ID          string                       `json:"id"`
	RequesterID string                       `json:"requester_id"`
	RecipientID string                       `json:"recipient_id"`
	Status      sharedEnums.StatusFriendship `json:"status"`
	Created_At  time.Time                    `json:"created_at"`
	Updated_At  time.Time                    `json:"updated_at"`
	EventType   constants.EventType          `json:"event_type"`
}
type BlockUserPayload struct {
	ID            string              `json:"id"`
	BlockerUserID string              `json:"blocker_user_id"`
	BlockedUserID string              `json:"blocked_user_id"`
	TypeBlock     string              `json:"type_block"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	EventType     constants.EventType `json:"event_type"`
}
