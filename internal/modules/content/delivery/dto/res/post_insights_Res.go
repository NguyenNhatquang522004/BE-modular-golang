package res

import (
	"time"
)

// PostInsightRes dùng để trả dữ liệu về cho client, ánh xạ 100% từ Entity
type PostInsightResv1 struct {
	PostID         string    `json:"post_id"`
	Reach          int64     `json:"reach"`
	Impressions    int64     `json:"impressions"`
	EngagementRate float32   `json:"engagement_rate"`
	ReactionsTotal int       `json:"reactions_total"`
	CommentsTotal  int       `json:"comments_total"`
	SharesTotal    int       `json:"shares_total"`
	ClicksTotal    int       `json:"clicks_total"`
	VideoViews3s   int       `json:"video_views_3s"`
	UpdatedAt      time.Time `json:"updated_at"`
}
