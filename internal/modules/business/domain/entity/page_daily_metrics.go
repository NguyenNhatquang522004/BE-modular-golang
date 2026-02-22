package entity

import (
	"time"

	"github.com/gocql/gocql"
)
const (
	CollectionPageDailyMetrics = "page_daily_metrics"
)
// PageDailyMetric đại diện cho bảng 'page_daily_metrics' trong Cassandra.
// Bảng này lưu trữ snapshot số liệu thống kê theo ngày (Time Series Data).
type PageDailyMetric struct {
	// =========================================================================
	// 1. PRIMARY KEY
	// =========================================================================

	// PARTITION KEY
	// ID của Page (Lưu ý: SQL là UUID, nếu dùng Mongo ObjectId thì cần đổi SQL thành TEXT)
	ID gocql.UUID `cql:"page_id" json:"page_id"`

	// CLUSTERING KEY
	// Ngày thống kê (VD: 2024-01-01). Sắp xếp giảm dần (DESC) để lấy ngày mới nhất.
	// Cassandra 'DATE' map về Go 'time.Time' (phần giờ phút giây sẽ được driver bỏ qua hoặc set về 00:00:00 UTC)
	MetricDate time.Time `cql:"metric_date" json:"metric_date"`

	// =========================================================================
	// 2. REACH METRICS (TIẾP CẬN)
	// =========================================================================
	// Dùng int64 cho BIGINT để tránh tràn số (Reach có thể lên tới hàng triệu)

	ReachTotal   int64 `cql:"reach_total" json:"reach_total"`     // Tổng tiếp cận
	ReachPaid    int64 `cql:"reach_paid" json:"reach_paid"`       // Quảng cáo
	ReachOrganic int64 `cql:"reach_organic" json:"reach_organic"` // Tự nhiên

	// =========================================================================
	// 3. ENGAGEMENT METRICS (TƯƠNG TÁC)
	// =========================================================================

	ImpressionsTotal int64 `cql:"impressions_total" json:"impressions_total"` // Số lần hiển thị

	// Dùng int cho INT (Đủ cho các chỉ số này trong 1 ngày)
	NewFollowers int `cql:"new_followers" json:"new_followers"`
	Unfollows    int `cql:"unfollows" json:"unfollows"`

	// =========================================================================
	// 4. CONVERSION METRICS (CHUYỂN ĐỔI)
	// =========================================================================

	ProfileViews  int `cql:"profile_views" json:"profile_views"`
	WebsiteClicks int `cql:"website_clicks" json:"website_clicks"`
	CTAClicks     int `cql:"cta_clicks" json:"cta_clicks"` // Click nút "Gửi tin nhắn", "Gọi ngay"...
}

// TableName trả về tên bảng trong Cassandra
func (PageDailyMetric) TableName() string {
	return CollectionPageDailyMetrics
}