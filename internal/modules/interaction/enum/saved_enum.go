package enum

//go:generate enumer -type=SavedTargetType -json -transform=snake -trimprefix=SavedTarget
type SavedTargetType int

const (
	SavedTargetPost  SavedTargetType = iota // 'post'
	SavedTargetReel                         // 'reel'
	SavedTargetVideo                        // 'video'
	// Có thể mở rộng thêm: SavedTargetProduct, SavedTargetComment...
)