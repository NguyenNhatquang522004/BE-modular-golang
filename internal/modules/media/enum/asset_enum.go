package enum

// =============================================================================
// 1. ASSET TYPE
// =============================================================================

//go:generate enumer -type=AssetType -json -transform=snake -trimprefix=Asset
type AssetType int

const (
	AssetImage   AssetType = iota // 'image'
	AssetVideo                    // 'video'
	AssetGif                      // 'gif'
	AssetSticker                  // 'sticker'
)

// =============================================================================
// 2. TAG STATUS (Trạng thái gắn thẻ bạn bè)
// =============================================================================

//go:generate enumer -type=TagStatus -json -transform=snake -trimprefix=TagStatus
type TagStatus int

const (
	TagStatusPending  TagStatus = iota // 'pending' (Chờ duyệt)
	TagStatusApproved                  // 'approved' (Đã hiện lên tường nhà người được tag)
)