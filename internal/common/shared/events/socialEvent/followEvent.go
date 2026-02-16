package socialEvent

type FollowCreatePayload struct {
	Follower_UserID string `json:"follower_id"`
	Followed_UserID string `json:"followed_id"`
	IsMuted         bool   `json:"is_muted"`
}
type FollowDeletePayload struct {
	Follower_UserID string `json:"follower_id"`
	Followed_UserID string `json:"followed_id"`
}
type FollowUpdatePayload struct {
	Follower_UserID string `json:"follower_id"`
	Followed_UserID string `json:"followed_id"`
	IsMuted         bool   `json:"is_muted"`
}
