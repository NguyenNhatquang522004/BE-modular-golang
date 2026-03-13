package sharedEnums
//go:generate enumer -type=CallType -json -transform=snake -trimprefix=CallType
type CallType int

const (
	CallTypeVoice CallType = iota // 'voice'
	CallTypeVideo                 // 'video'
)
