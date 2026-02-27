package sharedEnums

// =============================================================================
// 1. TARGET TYPE (Post / Comment)
// =============================================================================

//go:generate enumer -type=ReactionTarget -json -transform=snake -trimprefix=ReactionTarget
type ReactionTarget int

const (
	ReactionTargetPost    ReactionTarget = iota // 'post'
	ReactionTargetComment                       // 'comment'
	ReactionTargetAlbum                         // 'album'
	ReactionTargetStory                         // 'story'
	ReactionTargetReel                          // 'reel'
	ReactionTargetUnknown                       // 'unknown
)

// =============================================================================
// 2. REACTION CODE (Like, Love, ...)
// =============================================================================

//go:generate enumer -type=ReactionCode -json -transform=snake -trimprefix=Reaction
type ReactionCode int

const (
	ReactionLike    ReactionCode = iota // 'like'
	ReactionLove                        // 'love'
	ReactionHaha                        // 'haha'
	ReactionSad                         // 'sad'
	ReactionAngry                       // 'angry'
	ReactionWow                         // 'wow'
	ReactionUnknown                     // 'unknown'
)
