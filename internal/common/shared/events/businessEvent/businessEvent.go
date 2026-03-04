package businessEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type StatsPagePayload struct {
	UserID         string    `json:"user_id"`
	PageID         string    `json:"page_id"`
	FollowersCount int       `json:"followers_count"`
	LikesCount     int       `json:"likes_count"`
	RatingScore    float64   `json:"rating_score"`
	ReviewCount    int       `json:"review_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type FollowerPagePayload struct {
	ID       string                   `json:"id,omitempty" validate:"omitempty"`
	PageID   string                   `json:"page_id" validate:"required"`
	UserID   string                   `json:"user_id" validate:"required,uuid"` // Đảm bảo là UUID theo chuẩn Postgres
	Settings *FollowerSettingsPayload `json:"settings,omitempty"`
}
type FollowerSettingsPayload struct {
	NotificationLevel sharedEnums.NotificationLevel `json:"notification_level" validate:"required"`
	IsFavorite        bool                          `json:"is_favorite"`
}

type TopicDailyMetricsPagePayload struct {
	ID         string    `json:"page_id" binding:"required,uuid"`
	MetricDate time.Time `json:"metric_date" binding:"required"`

	ReachTotal   int64 `json:"reach_total" binding:"gte=0"`
	ReachPaid    int64 `json:"reach_paid" binding:"gte=0"`
	ReachOrganic int64 `json:"reach_organic" binding:"gte=0"`

	ImpressionsTotal int64 `json:"impressions_total" binding:"gte=0"`
	NewFollowers     int   `json:"new_followers" binding:"gte=0"`
	Unfollows        int   `json:"unfollows" binding:"gte=0"`

	ProfileViews  int `json:"profile_views" binding:"gte=0"`
	WebsiteClicks int `json:"website_clicks" binding:"gte=0"`
	CTAClicks     int `json:"cta_clicks" binding:"gte=0"`
}
