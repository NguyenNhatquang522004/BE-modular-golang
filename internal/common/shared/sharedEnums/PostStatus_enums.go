package sharedEnums

//go:generate enumer -type=PostStatus -json -transform=snake -trimprefix=Status
type PostStatus int

const (
	StatusPublished PostStatus = iota // 'published'
	StatusDraft                       // 'draft'
	StatusArchived                    // 'archived'
	StatusHidden                      // 'hidden'
)
