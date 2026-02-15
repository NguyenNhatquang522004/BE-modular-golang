package res

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"

// UserRecommendation: Kết quả gợi ý kết bạn (PYMK)
type UserRecommendation struct {
	User          *entity.UserNode `json:"user"`
	MutualFriends int              `json:"mutual_friends"`
	CommonGroups  int              `json:"common_groups"`
	MatchScore    float64          `json:"match_score"` // AI Prediction Score
	Reasons       []string         `json:"reasons"`     // ["Gần bạn", "Bạn bè của A"]
}

// TrendingTopic: Kết quả Trending
type TrendingTopic struct {
	Topic      *entity.TopicNode `json:"topic"`
	PostCount  int               `json:"post_count"`
	GrowthRate float64           `json:"growth_rate"`
}
