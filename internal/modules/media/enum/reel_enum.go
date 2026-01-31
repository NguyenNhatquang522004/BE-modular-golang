package enum

// =============================================================================
// 1. PROCESSING STATUS (Quan trọng cho Worker xử lý video)
// =============================================================================

//go:generate enumer -type=ProcessingStatus -json -transform=snake -trimprefix=Processing
type ProcessingStatus int

const (
	ProcessingPending    ProcessingStatus = iota // 'pending' (Mới upload, chưa xử lý)
	ProcessingProcessing                         // 'processing' (Đang transcode/nén)
	ProcessingActive                             // 'active' (Đã xong, user có thể xem)
	ProcessingFailed                             // 'failed' (Lỗi file)
)

// =============================================================================
// 2. REMIX TYPE
// =============================================================================

//go:generate enumer -type=RemixType -json -transform=snake -trimprefix=Remix
type RemixType int

const (
	RemixDuet  RemixType = iota // 'duet' (Chia đôi màn hình)
	RemixRemix                  // 'remix' (Lồng ghép)
)

// =============================================================================
// 3. PRIVACY (Tái sử dụng hoặc định nghĩa mới)
// =============================================================================

//go:generate enumer -type=ReelPrivacy -json -transform=snake -trimprefix=Privacy
type ReelPrivacy int

const (
	PrivacyPublic  ReelPrivacy = iota // 'public'
	PrivacyFriends                    // 'friends'
	PrivacyOnlyMe                     // 'only_me'
)