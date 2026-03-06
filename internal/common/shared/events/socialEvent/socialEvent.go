package socialEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"

type DeleteSocialRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}

type FollowerUserPayload struct {
	ID             string              `json:"id"`
	FollowerUserID string              `json:"follower_user_id"`
	FollowedUserID string              `json:"followed_user_id"`
	IsMuted        bool                `json:"is_muted"`
	EventType      constants.EventType `json:"event_type"`
}
