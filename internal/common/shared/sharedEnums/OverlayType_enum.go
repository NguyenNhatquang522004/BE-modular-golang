package sharedEnums

//go:generate enumer -type=OverlayType -json -transform=snake -trimprefix=Overlay
type OverlayType int

const (
	OverlayText     OverlayType = iota // 'text'
	OverlayMusic                       // 'music'
	OverlayMention                     // 'mention'
	OverlayLocation                    // 'location'
	OverlayPoll                        // 'poll'
	OverlayQuestion                    // 'question'
)
