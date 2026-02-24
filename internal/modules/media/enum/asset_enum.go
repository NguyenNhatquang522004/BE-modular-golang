package enum

// =============================================================================
// 2. TAG STATUS (Trạng thái gắn thẻ bạn bè)
// =============================================================================

//go:generate enumer -type=TagStatus -json -transform=snake -trimprefix=TagStatus
type TagStatus int

const (
	TagStatusPending  TagStatus = iota // 'pending' (Chờ duyệt)
	TagStatusApproved                  // 'approved' (Đã hiện lên tường nhà người được tag)
)