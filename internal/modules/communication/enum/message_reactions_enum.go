package enum

//go:generate enumer -type=MessageReactionCode -json -transform=snake -trimprefix=Reaction
type MessageReactionCode int

const (
	ReactionHeart   MessageReactionCode = iota // 'heart' (❤️)
	ReactionHaha                       // 'haha' (😆)
	ReactionSad                        // 'sad' (😢)
	ReactionWow                        // 'wow' (😮)
	ReactionAngry                      // 'angry' (😡)
	ReactionLike                       // 'like' (👍)
	ReactionDislike                    // 'dislike' (👎)
)