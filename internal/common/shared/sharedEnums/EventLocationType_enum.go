package sharedEnums
//go:generate enumer -type=EventLocationType -json -transform=snake -trimprefix=EventLocation
type EventLocationType int

const (
	EventLocationOnline  EventLocationType = iota // 'online'
	EventLocationOffline                      // 'offline'
)