package sharedEnums

//go:generate enumer -type=ConversationType -json -transform=snake -trimprefix=Type
type ConversationType int

const (
	TypePrivate ConversationType = iota // 'private' (1-1)
	TypeGroup                           // 'group' (Nhiều người)
	TypeChannel                         // 'channel' (Kênh trong nhóm lớn)
)