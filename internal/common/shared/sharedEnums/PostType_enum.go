package sharedEnums

//go:generate enumer -type=PostType -json -transform=snake -trimprefix=PostType
type PostType int

const (
	PostTypeText       PostType = iota // 'text'
	PostTypeMedia                      // 'media'
	PostTypeShare                      // 'share'
	PostTypeBackground                 // 'background'
	PostTypeQnA                        // 'qna'
	PostTypeLive                       // 'live'
)
