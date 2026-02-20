package req

import "fmt"

type PostInsightsReq struct {
	PostID string `json:"post_id" binding:"required,uuid4"` // UUID v4 string
}

type UpdatePostInsightsInteractionReq struct {
	*PostInsightsReq
	// Loại tương tác: "reaction", "comment", "share", "click", "video_view_3s"
	InteractionType string `json:"interaction_type" binding:"required,oneof=reaction comment share click video_view_3s"`
	// Số lượng tương tác mới (có thể là +1 hoặc -1 hoặc số cụ thể)
	Count int `json:"count" binding:"required"`
}
type UpdatePostInsightsLifeTimeReq struct {
	*PostInsightsReq
	// Các chỉ số trọn đời có thể cập nhật: "reach", "impressions", "engagement_rate"
	MetricType string `json:"metric_type" binding:"required,oneof=reach impressions"`
	// Giá trị mới của chỉ số (có thể là số cụ thể hoặc tăng/giảm)
	Value float64 `json:"value" binding:"required"`
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
