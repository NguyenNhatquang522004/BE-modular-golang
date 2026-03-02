package enum

// =============================================================================
// 1. PRIVACY LEVEL
// =============================================================================

//go:generate enumer -type=GroupPrivacy -json -transform=snake -trimprefix=Privacy
type GroupPrivacy int

const (
	PrivacyPublic  GroupPrivacy = iota // 'public' (Ai cũng tìm thấy và xem bài)
	PrivacyPrivate                     // 'private' (Tìm thấy, nhưng phải join mới xem bài)
	PrivacySecret                      // 'secret' (Không tìm thấy, chỉ invite)
)
