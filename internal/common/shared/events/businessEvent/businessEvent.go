package businessEvent

import "time"

type StatsPagePayload struct {
	UserID         string    `json:"user_id"`
	PageID         string    `json:"page_id"`
	FollowersCount int       `json:"followers_count"`
	LikesCount     int       `json:"likes_count"`
	RatingScore    float64   `json:"rating_score"`
	ReviewCount    int       `json:"review_count"`
	CreatedAt      time.Time `json:"created_at"`
}
