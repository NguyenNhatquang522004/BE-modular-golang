package sharedEnums

//go:generate enumer -type=ActivityType -json -transform=snake -trimprefix=Activity
type ActivityType int

const (
	ActivityFeeling   ActivityType = iota // 'feeling' (đang cảm thấy)
	ActivityWatching                      // 'watching' (đang xem)
	ActivityReading                       // 'reading' (đang đọc)
	ActivityListening                     // 'listening' (đang nghe)
	ActivityTraveling                     // 'traveling' (đang đi tới)
	ActivityEating                        // 'eating' (đang ăn)
	ActivityCelebrating                   // 'celebrating' (đang chúc mừng)
)