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

// =============================================================================
// 2. APPROVER PERMISSION (Ai có quyền duyệt thành viên?)
// =============================================================================

//go:generate enumer -type=GroupApprover -json -transform=snake -trimprefix=Approver
type GroupApprover int

const (
	ApproverAdmin    GroupApprover = iota // 'admin' (Chỉ Admin/Mod)
	ApproverEveryone                      // 'everyone' (Thành viên cũng được duyệt - Ít dùng nhưng Facebook có)
)