package enum

// =============================================================================
// 1. STORY MEDIA TYPE
// =============================================================================

//go:generate enumer -type=StoryMediaType -json -transform=snake -trimprefix=StoryMedia
type StoryMediaType int

const (
	StoryMediaImage StoryMediaType = iota // 'image'
	StoryMediaVideo                       // 'video'
)

// =============================================================================
// 2. PRIVACY TYPE
// =============================================================================

//go:generate enumer -type=StoryPrivacyType -json -transform=snake -trimprefix=StoryPrivacy
type StoryPrivacyType int

const (
	StoryPrivacyPublic       StoryPrivacyType = iota // 'public'
	StoryPrivacyFriends                              // 'friends'
	StoryPrivacyCloseFriends                         // 'close_friends'
	StoryPrivacyCustom                               // 'custom'
)

// =============================================================================
// 3. OVERLAY TYPE (Stickers, Polls...)
// =============================================================================

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