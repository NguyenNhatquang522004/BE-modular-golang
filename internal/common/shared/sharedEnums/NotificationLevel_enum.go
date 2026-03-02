package sharedEnums

//go:generate enumer -type=NotificationLevel -json -transform=snake -trimprefix=Notif
type NotificationLevel int

const (
	NotifAll       NotificationLevel = iota // 'all' (Tất cả bài viết)
	NotifHighlight                          // 'highlight' (Chỉ bài nổi bật - Mặc định)
	NotifOff                                // 'off' (Tắt thông báo)
)
