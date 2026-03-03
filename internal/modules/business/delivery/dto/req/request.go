package req

import "time"

type CreateAdAccountRequest struct {
	*PageReq
}
type UpdateAdAccountRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*PageReq
}
type DeletePageRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	PageID       string `json:"page_id" validate:"required"`
}
type StatsPageRequest struct {
	UserID         string    `json:"user_id"`
	PageID         string    `json:"page_id"`
	FollowersCount int       `json:"followers_count"`
	LikesCount     int       `json:"likes_count"`
	RatingScore    float64   `json:"rating_score"`
	ReviewCount    int       `json:"review_count"`
	CreatedAt      time.Time `json:"created_at"`
}
