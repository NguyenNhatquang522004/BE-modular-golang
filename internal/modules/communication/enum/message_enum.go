package enum

//go:generate enumer -type=MessageType -json -transform=snake -trimprefix=Msg
type MessageType int

const (
	MsgText       MessageType = iota // 'text'
	MsgImage                         // 'image'
	MsgVideo                         // 'video'
	MsgAudio                         // 'audio'
	MsgFile                          // 'file'
	MsgLocation                      // 'location'
	MsgStoryReply                    // 'story_reply'
	MsgSticker                       // 'sticker'
)
