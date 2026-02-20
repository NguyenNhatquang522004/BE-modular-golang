package mapper

import "fmt"

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
