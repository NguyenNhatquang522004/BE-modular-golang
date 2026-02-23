package contentEvent

type PostDeletePayload struct {
	PostID string `json:"post_id"`
}
type PostStatsPayload struct {
	PostID      string `json:"post_id"`
	TotalLikes  int    `json:"total_likes"`
	TotalShares int    `json:"total_shares"`
	TotalViews  int    `json:"total_views"`
}
