package graphEvent

type ScoredPost struct {
	PostID string  `json:"post_id"`
	Score  float64 `json:"score"`
}
type SuggestedUser struct {
	UserID             string  `json:"user_id"`
	MutualFriendsCount int     `json:"mutual_friends_count"`
	Score              float64 `json:"score"`
}
type FeedItem struct {
	Type string `json:"type"` // "post" hoặc "pymk"
	Data any    `json:"data"` // Chứa entity.ScoredPost HOẶC mảng entity.SuggestedUser
}

type NewsFeedResponse struct {
	Items      []FeedItem `json:"items"`
	NextCursor string     `json:"next_cursor,omitempty"`
}
type PostMeta struct {
	AuthorID  string
	Topic     string
	RiskScore float64
	CreatedAt int64
}
type CandidateSignals struct {
	S1_Score float64
	S2_Score float64
	S3_Score float64
	S4_Score float64
	S5_Score float64
	S6_Score float64
	S7_Score float64
	
	// Số lượng bạn chung thực tế (Lấy từ S1 để ưu tiên hiển thị UI)
	RealMutualFriends int 
	
	// Số lượng chiến lược (Strategies) mà ứng viên này xuất hiện
	// Dùng để buff điểm chéo (Cross-Signal Boost)
	HitCount int 
}