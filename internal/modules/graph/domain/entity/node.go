package entity

type UserNode struct {
	UserID       string    `json:"user_id"`        // UUID
	CreatedAt    int64     `json:"created_at"`     // Timestamp
	LastActiveAt int64     `json:"last_active_at"` // Timestamp
	IsVerified   bool      `json:"is_verified"`
	Bio          string    `json:"bio,omitempty"` // Thông tin giới thiệu bản thân
	Embedding    []float32 `json:"embedding"`     // Vector 128D
	// Analytics Fields (Tính toán từ GDS)
	PageRankScore float64 `json:"page_rank_score,omitempty"` // Độ uy tín
	CommunityID   int64   `json:"community_id,omitempty"`    // ID cụm cộng đồng
}

// TopicNode (Interest Graph)
type TopicNode struct {
	TopicID       string  `json:"topic_id"`
	Name          string  `json:"name"`           // Unique (Hashtag)
	TrendingScore float64 `json:"trending_score"` // Độ hot hiện tại
	// Note: Category cũ đã bỏ, thay bằng quan hệ CHILD_OF bên dưới
}

// PostNode (Short-term / Viral Content - [MỚI])
// Chỉ lưu Node này khi bài viết đang Trending
type PostNode struct {
	PostID    string `json:"post_id"`
	CreatedAt int64  `json:"created_at"`
	TTL       int64  `json:"ttl"` // Thời điểm tự hủy (Time-To-Live)
}

// GroupNode (Community)
type GroupNode struct {
	GroupID     string `json:"group_id"`
	Privacy     string `json:"privacy"` // public, closed, secret
	MemberCount int    `json:"member_count"`
}

// PageNode (Brand/Fanpage)
type PageNode struct {
	PageID     string  `json:"page_id"`
	CategoryID string  `json:"category_id"`
	Rating     float64 `json:"rating"`
}

// PhoneContactNode (Identity)
type PhoneContactNode struct {
	PhoneHash string `json:"phone_hash"` // SHA256
}

// LocationNode (Geo)
type LocationNode struct {
	CityID      string `json:"city_id"`
	CountryCode string `json:"country_code"`
	GeoHash     string `json:"geo_hash"`
}
