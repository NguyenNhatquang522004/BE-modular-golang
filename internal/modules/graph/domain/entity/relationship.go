package entity

// ==========================================
// 2. RELATIONSHIP ENTITIES (Thực thể Cạnh)
// ==========================================

// --- NHÓM XÃ HỘI & TĂNG TRƯỞNG ---

// FriendRel (Social)
type FriendRel struct {
	Since                int64   `json:"since"`
	Type                 string  `json:"type"` // normal, close_friend, family
	InteractionFrequency float64 `json:"interaction_frequency"`
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

// ChildOfRel (Topic Hierarchy) - [MỚI]
// VD: (Golang)-[:CHILD_OF {weight: 0.9}]->(Backend)
type ChildOfRel struct {
	Weight float64 `json:"weight"`
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

// UsedDeviceRel (User -> Device)
type UsedDeviceRel struct {
	LastUsedAt int64 `json:"last_used_at"`
	LoginCount int   `json:"login_count"`
}

// HasContactRel (User -> PhoneContact) - [MỚI]
type HasContactRel struct {
	UploadedAt int64 `json:"uploaded_at"`
}

// --- NHÓM TÍN HIỆU TIÊU CỰC ---

type BlockRel struct {
	Since int64 `json:"since"`
}

type MuteRel struct {
	Since int64 `json:"since"`
}

type HiddenPostRel struct {
	Count        int   `json:"count"`
	LastHiddenAt int64 `json:"last_hidden_at"`
}

type ReportedRel struct {
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}
