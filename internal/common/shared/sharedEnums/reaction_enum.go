package sharedEnums

// =============================================================================
// 1. TARGET TYPE (Post / Comment)
// =============================================================================

//go:generate enumer -type=ReactionTarget -json -transform=snake -trimprefix=ReactionTarget
type ReactionTarget int

const (
	ReactionTargetPost         ReactionTarget = iota // 'post'
	ReactionTargetComment                            // 'comment'
	ReactionTargetAlbum                              // 'album'
	ReactionTargetStory                              // 'story'
	ReactionTargetViewStory                          // 'view_story'
	ReactionTargetReplyStory                         // 'Reply_story'
	ReactionTargetReel                               // 'reel'
	ReactionTargetViewReel                           // 'view_reel'
	ReactionTargetShareReel                          // 'share_reel'
	ReactionTargetCommentReel                        // 'comment_reel'
	ReactionTargetSaveReel                           // 'save_reel'
	ReactionTargetUnknown                            // 'unknown'
	ReactionTargetLive                               // 'live'
	ReactionTargetViewLive                           // 'view_live'
	ReactionTargetCommentLive                        // 'comment_live'
	ReactionTargetFollowPage                         // 'follow_page'
	ReactionTargetUnFollowPage                       // 'unfollow_page'
	ReactionTargetSharePost                          // 'share_post'
	ReactionTargetViewPost                           // 'view_post'
	ReactionTargetReplyComment                       // 'reply_comment'
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
