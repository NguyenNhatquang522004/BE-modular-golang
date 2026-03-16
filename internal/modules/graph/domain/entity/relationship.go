package entity

// ==========================================
// 2. RELATIONSHIP ENTITIES (Thực thể Cạnh)
// ==========================================
// AuthoredRel (User -> Post)
type AuthoredRel struct {
	PostID    string `json:"post_id"`
	CreatedAt int64  `json:"created_at"`
}

// PublishedRel (Page -> Post) - [BỔ SUNG ĐỂ PAGE ĐĂNG BÀI]
type PublishedRel struct {
	CreatedAt int64 `json:"created_at"`
}

// HasTopicRel (Post -> Topic) - [ĐÃ SỬA LỖI]
type HasTopicRel struct {
	ConfidenceScore float64 `json:"confidence_score"` // Model AI tự tin bao nhiêu %
}

type PostedInRel struct {
	GroupID string `json:"group_id"`
	Privacy string `json:"privacy"` // public, closed, secret
}
type LocatedInRel struct {
	UpdatedAt int64 `json:"updated_at"`
}

// --- NHÓM XÃ HỘI & TĂNG TRƯỞNG ---

// FriendRel (Social)
type FriendRel struct {
	Since                int64   `json:"since"`
	Type                 string  `json:"type"` // normal, close_friend, family
	InteractionFrequency float64 `json:"interaction_frequency"`
}
type BlockRel struct {
	Since int64 `json:"since"`
}

// FollowRel (Social)
type FollowRel struct {
	Since  int64  `json:"since"`
	Source string `json:"source"` // profile, search, suggested
}

// InvitedRel (Growth & Security) - [MỚI QUAN TRỌNG]
type InvitedRel struct {
	Timestamp int64  `json:"timestamp"`
	Code      string `json:"code"` // Mã mời hoặc method (email/sms)
}

// --- NHÓM TƯƠNG TÁC (RANKING CORE) ---

// InteractedWithRel (Long-term Memory / Aggregated)
type InteractedWithRel struct {
	LastInteractionAt int64   `json:"last_interaction_at"`
	LikeCount         int     `json:"like_count"`
	CommentCount      int     `json:"comment_count"`
	MessageCount      int     `json:"message_count"`
	ShareCount        int     `json:"share_count"`
	ProfileViewCount  int     `json:"profile_view_count"`
	AffinityScore     float64 `json:"affinity_score"` // Ranking Score
}

// InteractedRecentlyRel (Short-term Memory / Real-time) - [MỚI]
type InteractedRecentlyRel struct {
	Timestamp int64   `json:"timestamp"`
	Type      string  `json:"type"`   // view, like, share
	Weight    float64 `json:"weight"` // Trọng số tức thời
}

// --- NHÓM SỞ THÍCH & NỘI DUNG ---

// InterestedInRel (User -> Topic)
type InterestedInRel struct {
	Score         float64 `json:"score"` // 0.0 - 1.0
	LastEngagedAt int64   `json:"last_engaged_at"`
}

type RelatedToRel struct {
    SimilarityScore float64 `json:"similarity_score"` // Cosine Similarity (VD: 0.85)
}
// MemberOfRel (User -> Group)
type MemberOfRel struct {
	Role     string `json:"role"` // admin, member
	JoinedAt int64  `json:"joined_at"`
}

// LikesPageRel (User -> Page) - [MỚI]
type LikesPageRel struct {
	Since int64 `json:"since"`
}

// --- NHÓM ĐỊNH DANH & BẢO MẬT ---

// HasContactRel (User -> PhoneContact) - [MỚI]
type HasContactRel struct {
	UploadedAt int64 `json:"uploaded_at"`
}
