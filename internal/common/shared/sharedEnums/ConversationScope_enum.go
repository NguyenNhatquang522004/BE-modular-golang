package sharedEnums

//go:generate enumer -type=ConversationScope -json -transform=snake -trimprefix=Scope
type ConversationScope int

const (
	ScopeMessenger        ConversationScope = iota // 'messenger' (Chat thường)
	ScopeCommunityChannel                          // 'community_channel' (Kênh trong nhóm lớn)
)
