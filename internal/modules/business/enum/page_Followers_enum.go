package enum

// =============================================================================
// 1. FOLLOW TYPE
// =============================================================================
// Facebook phân tách: Like (Thích) và Follow (Theo dõi).
// - Like: Thường tự động bật Follow.
// - Follow: Có thể Follow mà không Like (để cập nhật tin tức).

//go:generate enumer -type=FollowType -json -transform=snake -trimprefix=FollowType
type FollowType int

const (
	FollowTypeFollow FollowType = iota // 'follow'
	FollowTypeLike                     // 'like'
)

// =============================================================================
// 2. NOTIFICATION LEVEL
// =============================================================================
// Mức độ nhận thông báo từ Page
