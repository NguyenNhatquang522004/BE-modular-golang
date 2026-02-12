package entity



// ==========================================
// 3. COMPOSITE MODELS (Dữ liệu trả về API)
// ==========================================

// UserRecommendation: Kết quả gợi ý kết bạn (PYMK)
type UserRecommendation struct {
	User          UserNode `json:"user"`
	MutualFriends int      `json:"mutual_friends"`
	CommonGroups  int      `json:"common_groups"`
	MatchScore    float64  `json:"match_score"` // AI Prediction Score
	Reasons       []string `json:"reasons"`     // ["Gần bạn", "Bạn bè của A"]
}

// TrendingTopic: Kết quả Trending
type TrendingTopic struct {
	Topic      TopicNode `json:"topic"`
	PostCount  int       `json:"post_count"`
	GrowthRate float64   `json:"growth_rate"`
}