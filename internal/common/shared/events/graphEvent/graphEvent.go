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
