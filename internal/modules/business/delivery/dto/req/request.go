package req

import "time"

type CreatePageRequest struct {
	*PageReq
}
type UpdatePageRequest struct {
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
type CreatePageRoleRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*PageRoleReq
}
type UpdatePageRoleRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*PageRoleReq
}
type DeletePageRoleRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	PageID       string `json:"page_id" validate:"required"`
	UserID       string `json:"user_id" validate:"required"`
}

type CreatePageFollowerRequest struct {
	*PageFollowerReq
}

type MetricRequest struct {
}

type CreateAdAccountRequest struct {
	*AdAccountReq
}

type UpdateAdAccountRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*AdAccountReq
}
type UpdateBalanceRequest struct {
	UserActionID string  `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	AccountID    string  `json:"account_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
}
type CreateAdCampainRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*AdCampainReq
}
