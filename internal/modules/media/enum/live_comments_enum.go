package enum

//go:generate enumer -type=UserBadge -json -transform=snake -trimprefix=Badge
type UserBadge int

const (
	BadgeModerator UserBadge = iota // 'moderator' (Quản trị viên)
	BadgeTopFan                     // 'top_fan' (Fan cứng)
	BadgeGifter                     // 'gifter' (Người tặng quà nhiều)
	BadgeVerified                   // 'verified' (Tích xanh)
	BadgeStaff                      // 'staff' (Nhân viên hệ thống)
)