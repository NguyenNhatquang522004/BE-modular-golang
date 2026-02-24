package sharedEnums

//go:generate enumer -type=MediaType -json -transform=snake -trimprefix=MediaType
type MediaType int

const (
	MediaTypeImage   MediaType = iota // 'image'
	MediaTypeVideo                    // 'video'
	MediaTypeGif                      // 'gif'
	MediaTypeSticker                  // 'sticker'

)
