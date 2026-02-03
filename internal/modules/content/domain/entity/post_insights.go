package entity

import (
	"time"

	"github.com/gocql/gocql"
)
const (
	collectionnamepostinsight = "PostInsight"
)
// PostInsight đại diện cho bảng 'post_insights' trong Cassandra
// Dùng để lưu trữ số liệu phân tích hiệu năng bài viết (High Write/Read throughput)
type PostInsight struct {
	// 1. PARTITION KEY
	// Sử dụng gocql.UUID để tương thích tốt nhất với driver Cassandra
	PostID gocql.UUID `cql:"post_id" json:"post_id"`

	// 2. LIFETIME SNAPSHOT (Chỉ số trọn đời)
	// BIGINT trong Cassandra tương ứng int64 trong Go
	Reach          int64   `cql:"reach" json:"reach"`
	Impressions    int64   `cql:"impressions" json:"impressions"`
	EngagementRate float32 `cql:"engagement_rate" json:"engagement_rate"` // FLOAT -> float32

	// 3. INTERACTION DETAILS (Chi tiết tương tác)
	// INT trong Cassandra tương ứng int trong Go
	ReactionsTotal int `cql:"reactions_total" json:"reactions_total"`
	CommentsTotal  int `cql:"comments_total" json:"comments_total"`
	SharesTotal    int `cql:"shares_total" json:"shares_total"`
	ClicksTotal    int `cql:"clicks_total" json:"clicks_total"`
	VideoViews3s   int `cql:"video_views_3s" json:"video_views_3s"`

	// 4. METADATA
	UpdatedAt time.Time `cql:"updated_at" json:"updated_at"`
}

func (PostSetting) Collectionnamepostinsight() string {
	return collectionnamepostinsight
}