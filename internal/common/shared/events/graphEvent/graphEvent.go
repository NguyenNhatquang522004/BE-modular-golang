package graphEvent

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

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
type TopicLinkInput struct {
	Name            string  `json:"name"`             // Tên Topic đã được chuẩn hóa (Canonical Name)
	ConfidenceScore float64 `json:"confidence_score"` // Độ tự tin AI
}
type PostNodePayload struct {
	UserID     string
	PostID     string                  `json:"post_id"`
	TargetType sharedEnums.ContextType `json:"target_type"`
	CreatedAt  int64                   `json:"created_at"`
	TTL        int64                   `json:"ttl"` // Thời điểm tự hủy (Time-To-Live)
	Topic      []TopicLinkInput        `json:"topic,omitempty"`
	EventType  constants.EventType     `json:"event_type"`
}

type DeletePostNodePayload struct {
	PostID   string `json:"post_id"`
	DeleteAt int64  `json:"delete_at"`
}

type UserNodePayload struct {
	UserID       string    `json:"user_id"`        // UUID
	CreatedAt    int64     `json:"created_at"`     // Timestamp
	LastActiveAt int64     `json:"last_active_at"` // Timestamp
	IsVerified   bool      `json:"is_verified"`
	Bio          string    `json:"bio,omitempty"` // Thông tin giới thiệu bản thân
	Embedding    []float32 `json:"embedding"`     // Vector 128D
	PhoneContact string    `json:"phone_contact,omitempty"`
	City         string    `json:"city,omitempty"`
	Country      string    `json:"country,omitempty"`
	GeoHash      string    `json:"geo_hash,omitempty"`
}

type DeleteUserNodePayload struct {
	UserID   string `json:"user_id"`
	DeleteAt int64  `json:"delete_at"`
}

type GroupNodePayload struct {
	GroupID     string `json:"group_id"`
	Privacy     string `json:"privacy"` // public, closed, secret
	MemberCount int    `json:"member_count"`
}

type DeleteGroupNodePayload struct {
	GroupID  string `json:"group_id"`
	DeleteAt int64  `json:"delete_at"`
}

type PageNodePayload struct {
	PageID     string  `json:"page_id"`
	CategoryID string  `json:"category_id"`
	Rating     float64 `json:"rating"`
}

type DeletePageNodePayload struct {
	PageID   string `json:"page_id"`
	DeleteAt int64  `json:"delete_at"`
}
type LocationNodePayload struct {
	CityID      string `json:"city_id"`
	CountryCode string `json:"country_code"`
	GeoHash     string `json:"geo_hash"`
}
type DeleteLocationNodePayload struct {
	CityID   string `json:"city_id"`
	DeleteAt int64  `json:"delete_at"`
}
type InteractionRecentPayload struct {
	UserID     string                     `json:"user_id"`
	PostID     string                     `json:"post_id"`
	TargetType sharedEnums.ContextType    `json:"target_type"`
	Timestamp  int64                      `json:"timestamp"`
	Type       sharedEnums.ReactionTarget `json:"type"`   // view, like, share
	Weight     float64                    `json:"weight"` // Trọng số tức thời
}
type InteractionPayload struct {
	UserID     string                  `json:"user_id"`
	TargetID   string                  `json:"target_id"`
	TargetType sharedEnums.ContextType `json:"target_type"`
	Like       int                     `json:"like"`
	Comment    int                     `json:"comment"`
	Share      int                     `json:"share"`
	Message    int                     `json:"message"`
	View       int                     `json:"view"`
	CreatedAt  int64                   `json:"created_at"`
	Flag       bool                    `json:"flag"` // true = cộng, false = trừ
}

type SharedPostPayload struct {
	UserID         string                  `json:"user_id"`
	NewSharePostID string                  `json:"new_share_post_id"`
	OriginalPostID string                  `json:"original_post_id"`
	ParentPostID   string                  `json:"parent_post_id"` // Có thể là post gốc hoặc post đã share trước đó
	CreatedAt      int64                   `json:"created_at"`
	TargetID       string                  `json:"target_id"`
	TargetType     sharedEnums.ContextType `json:"target_type"`
}
type DeleteSharedPostPayload struct {
	SharedPostID string `json:"shared_post_id"`
	DeleteAt     int64  `json:"delete_at"`
}

type FriendshipPayload struct {
	UserA          string              `json:"user_a"`
	UserB          string              `json:"user_b"`
	FriendshipType string              `json:"friendship_type"`
	Since          int64               `json:"since"`
	EventType      constants.EventType `json:"event_type"`
}

type DeleteFriendshipPayload struct {
	UserA    string `json:"user_a"`
	UserB    string `json:"user_b"`
	DeleteAt int64  `json:"delete_at"`
}

type BlockPayload struct {
	SourceUserID string                  `json:"source_user_id"`
	TargetUserID string                  `json:"target_user_id"`
	TargetType   sharedEnums.ContextType `json:"target_type"`
	Since        int64                   `json:"since"`
}

type FollowPayload struct {
	FollowerID string                  `json:"follower_id"`
	TargetID   string                  `json:"target_id"`
	TargetType sharedEnums.ContextType `json:"target_type"`
	Since      int64                   `json:"since"`
}
type JoinGroupPayload struct {
	UserID    string              `json:"user_id"`
	GroupID   string              `json:"group_id"`
	Role      string              `json:"role"` // member, admin
	JoinedAt  int64               `json:"joined_at"`
	EventType constants.EventType `json:"event_type"`
}

type DeleteJoinGroupPayload struct {
	UserID   string `json:"user_id"`
	GroupID  string `json:"group_id"`
	DeleteAt int64  `json:"delete_at"`
}
type LikePagePayload struct {
	UserID    string              `json:"user_id"`
	PageID    string              `json:"page_id"`
	Since     int64               `json:"since"`
	EventType constants.EventType `json:"event_type"`
}
type DeleteLikePagePayload struct {
	UserID   string `json:"user_id"`
	PageID   string `json:"page_id"`
	DeleteAt int64  `json:"delete_at"`
}
