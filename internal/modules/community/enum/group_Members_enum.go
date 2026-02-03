package enum

// =============================================================================
// 1. MEMBER ROLE
// =============================================================================

//go:generate enumer -type=MemberRole -json -transform=snake -trimprefix=Role
type MemberRole int

const (
	RoleMember    MemberRole = iota // 'member'
	RoleAdmin                       // 'admin'
	RoleModerator                   // 'moderator'
)

// =============================================================================
// 2. MEMBER STATUS
// =============================================================================

//go:generate enumer -type=MemberStatus -json -transform=snake -trimprefix=Status
type MemberStatus int

const (
	StatusActive  MemberStatus = iota // 'active'
	StatusPending                     // 'pending' (Xin vào, chờ duyệt)
	StatusInvited                     // 'invited' (Được mời, chờ đồng ý)
	StatusBanned                      // 'banned' (Bị chặn khỏi nhóm)
	StatusMuted                       // 'muted' (Bị cấm chat tạm thời)
)

// =============================================================================
// 3. MEMBER BADGE
// =============================================================================

//go:generate enumer -type=MemberBadge -json -transform=snake -trimprefix=Badge
type MemberBadge int

const (
	BadgeFoundingMember     MemberBadge = iota // 'founding_member'
	BadgeTopContributor                        // 'top_contributor'
	BadgeConversationStarter                   // 'conversation_starter'
	BadgeVisualStoryteller                     // 'visual_storyteller' (Người hay đăng ảnh đẹp)
	BadgeNewMember                             // 'new_member'
)