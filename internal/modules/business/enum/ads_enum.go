package enum

// =============================================================================
// AD STATUS (Quy trình xét duyệt)
// =============================================================================

//go:generate enumer -type=AdStatus -json -transform=snake -trimprefix=AdStatus
type AdStatus int

const (
	AdStatusReviewing AdStatus = iota // 'reviewing' (Mặc định khi mới tạo)
	AdStatusActive                    // 'active' (Đã duyệt và đang chạy)
	AdStatusRejected                  // 'rejected' (Bị từ chối do vi phạm policy)
	AdStatusPaused                    // 'paused' (Người dùng tự tắt)
)