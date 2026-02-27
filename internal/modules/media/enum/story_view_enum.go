package enum

//go:generate enumer -type=StoryInteractionType -json -transform=snake -trimprefix=Interaction
type StoryInteractionType int

const (
	InteractionView     StoryInteractionType = iota // 'view' (Chỉ xem)
	InteractionReaction                             // 'reaction' (Thả tim/haha)
	InteractionPollVote                             // 'poll_vote' (Bình chọn)
	InteractionReply                                // 'reply' (Trả lời)
)
