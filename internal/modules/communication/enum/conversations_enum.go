package enum

// =============================================================================
// 1. CONVERSATION TYPE & SCOPE
// =============================================================================

//go:generate enumer -type=ConversationType -json -transform=snake -trimprefix=Type
type ConversationType int

const (
	TypePrivate ConversationType = iota // 'private' (1-1)
	TypeGroup                           // 'group' (Nhiều người)
)

//go:generate enumer -type=ConversationScope -json -transform=snake -trimprefix=Scope
type ConversationScope int

const (
	ScopeMessenger        ConversationScope = iota // 'messenger' (Chat thường)
	ScopeCommunityChannel                          // 'community_channel' (Kênh trong nhóm lớn)
)

// =============================================================================
// 2. STATUS
// =============================================================================


