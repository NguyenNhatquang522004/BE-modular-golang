package sharedEnums

//go:generate enumer -type=MediaType -json -transform=snake -trimprefix=MediaType
type MediaType int

const (
	MediaTypeImage   MediaType = iota // 'image'
	MediaTypeVideo                    // 'video'
	MediaTypeGif                      // 'gif'
	MediaTypeSticker                  // 'sticker'
	MediaTypeAudio                    // 'audio'
	MediaTypeText                     // 'text'
	MediaTypePDF                      // 'pdf'
	MediaTypeDOCX                     // 'docx' (Word)
	MediaTypeXLSX                     // 'xlsx' (Excel)
	MediaTypePPTX                     // 'pptx' (PowerPoint)
	MediaTypeZIP                      // 'zip'  (Archive)
	MediaTypeRAR                      // 'rar'
	MediaTypeTXT                      // 'txt'
	MediaTypeCSV                      // 'csv'
)
