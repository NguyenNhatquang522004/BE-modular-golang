package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/gocql/gocql"
)

// ToEntity chuyển đổi từ Request DTO sang Entity để lưu vào Cassandra
func ToEntityPostInsight(r *req.PostInsightReqv1) (*entity.PostInsight, error) {
	if r == nil {
		return nil, fmt.Errorf("request payload is nil")
	}

	// Parse chuỗi UUID sang kiểu gocql.UUID
	postID, err := gocql.ParseUUID(r.PostID)
	if err != nil {
		return nil, fmt.Errorf("invalid post_id format: %w", err)
	}

	// Best practice: Tự động gán thời gian hiện tại nếu client không truyền UpdatedAt
	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	return &entity.PostInsight{
		PostID:         postID,
		Reach:          r.Reach,
		Impressions:    r.Impressions,
		EngagementRate: r.EngagementRate,
		ReactionsTotal: r.ReactionsTotal,
		CommentsTotal:  r.CommentsTotal,
		SharesTotal:    r.SharesTotal,
		ClicksTotal:    r.ClicksTotal,
		VideoViews3s:   r.VideoViews3s,
		UpdatedAt:      updatedAt,
	}, nil
}

// UpdateToEntity cập nhật dữ liệu từ Request DTO vào một Entity đang có sẵn
func UpdateToEntityPostInsight(r *req.PostInsightReqv1, e *entity.PostInsight) error {
	if r == nil || e == nil {
		return fmt.Errorf("request payload or entity is nil")
	}

	// Cập nhật ID nếu client có gửi lên ID mới (thường ít khi sửa khóa chính, nhưng cứ ánh xạ đủ 100% theo yêu cầu)
	if r.PostID != "" {
		postID, err := gocql.ParseUUID(r.PostID)
		if err != nil {
			return fmt.Errorf("invalid post_id format: %w", err)
		}
		e.PostID = postID
	}

	e.Reach = r.Reach
	e.Impressions = r.Impressions
	e.EngagementRate = r.EngagementRate
	e.ReactionsTotal = r.ReactionsTotal
	e.CommentsTotal = r.CommentsTotal
	e.SharesTotal = r.SharesTotal
	e.ClicksTotal = r.ClicksTotal
	e.VideoViews3s = r.VideoViews3s

	// Cập nhật UpdatedAt thành thời gian hiện tại mỗi khi có update
	e.UpdatedAt = time.Now()

	return nil
}

// ReqToRes chuyển đổi trực tiếp từ Request DTO sang Response DTO (Theo đúng yêu cầu của bạn)
func ReqToResPostInsight(r *req.PostInsightReqv1) *res.PostInsightResv1 {
	if r == nil {
		return nil
	}

	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	return &res.PostInsightResv1{
		PostID:         r.PostID,
		Reach:          r.Reach,
		Impressions:    r.Impressions,
		EngagementRate: r.EngagementRate,
		ReactionsTotal: r.ReactionsTotal,
		CommentsTotal:  r.CommentsTotal,
		SharesTotal:    r.SharesTotal,
		ClicksTotal:    r.ClicksTotal,
		VideoViews3s:   r.VideoViews3s,
		UpdatedAt:      updatedAt,
	}
}

// EntityToRes (Bonus Best Practice): Hàm chuyển từ Entity sang Response DTO (Rất hay dùng khi query từ DB lên)
func EntityToResPostInsight(e *entity.PostInsight) *res.PostInsightResv1 {
	if e == nil {
		return nil
	}
	return &res.PostInsightResv1{
		PostID:         e.PostID.String(),
		Reach:          e.Reach,
		Impressions:    e.Impressions,
		EngagementRate: e.EngagementRate,
		ReactionsTotal: e.ReactionsTotal,
		CommentsTotal:  e.CommentsTotal,
		SharesTotal:    e.SharesTotal,
		ClicksTotal:    e.ClicksTotal,
		VideoViews3s:   e.VideoViews3s,
		UpdatedAt:      e.UpdatedAt,
	}
}
func MapInteractionTypeToColumn(interactionType string) (string, error) {
	switch interactionType {
	case "reaction":
		return "reactions_total", nil
	case "comment":
		return "comments_total", nil
	case "share":
		return "shares_total", nil
	case "click":
		return "clicks_total", nil
	case "video_view_3s":
		return "video_views_3s", nil
	default:
		return "", fmt.Errorf("invalid interaction type: %s", interactionType)
	}
}

func MapMetricTypeToColumn(metricType string) (string, error) {
	switch metricType {
	case "reach":
		return "reach", nil
	case "impressions":
		return "impressions", nil
	case "engagement_rate":
		return "engagement_rate", nil
	default:
		return "", fmt.Errorf("invalid lifetime metric type: %s", metricType)
	}
}
