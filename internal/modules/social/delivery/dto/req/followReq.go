package req

type followerUseridRequest struct {
	FollowerUserID string `json:"follower_user_id" validate:"required,uuid4"`
}

type followedUserIDRequest struct {
	FollowedUserID string `json:"followed_user_id" validate:"required,uuid4"`
}
type FollowCreateRequest struct {
	followerUseridRequest
	followedUserIDRequest
}

type FollowDeleteSoftRequest struct {
	followerUseridRequest
}

type FollowDeleteHardRequest struct {
	followerUseridRequest
}

type FollowUpdateMuteRequest struct {
	followerUseridRequest
	followedUserIDRequest
	IsMuted bool `json:"is_muted" validate:"required"`
}

type FollowPaginationRequest struct {
	UserID string `form:"user_id" validate:"required,uuid4"`
	Cursor string `form:"cursor"`
	Limit  int    `form:"limit" validate:"gte=1,lte=100"`
}

