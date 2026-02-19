package entity

import (
	"fmt"
	"log"
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

func (PostInsight) Collectionnamepostinsight() string {
	return collectionnamepostinsight
}

func (p *PostInsight) EnsureTableExists(session *gocql.Session) error {
	// 1. Định nghĩa câu lệnh CQL
	// Sử dụng 'IF NOT EXISTS' để tránh lỗi nếu bảng đã có rồi.
	// Lưu ý: Cassandra yêu cầu xác định rõ Replication Strategy khi tạo KEYSPACE,
	// nhưng ở đây ta giả định Keyspace đã được config trong Session.
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			post_id         UUID,
			reach           bigint,
			impressions     bigint,
			engagement_rate float,
			reactions_total int,
			comments_total  int,
			shares_total    int,
			clicks_total    int,
			video_views_3s  int,
			updated_at      timestamp,
			PRIMARY KEY (post_id)
		) WITH compaction = { 'class' : 'LeveledCompactionStrategy' }
		AND comment = 'Auto-generated table for PostInsight';
	`, p.Collectionnamepostinsight())

	// 2. Thực thi query
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to auto-create table %s: %w", p.Collectionnamepostinsight(), err)
	}

	log.Printf("Successfully ensured table '%s' exists.", p.Collectionnamepostinsight())
	return nil
}
