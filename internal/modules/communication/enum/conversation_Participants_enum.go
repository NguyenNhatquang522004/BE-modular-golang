package enum

// =============================================================================
// PARTICIPANT ROLE
// =============================================================================

//go:generate enumer -type=ParticipantRole -json -transform=snake -trimprefix=Role
type ParticipantRole int

const (
	RoleMember    ParticipantRole = iota // 'member'
	RoleAdmin                            // 'admin' (Quản trị viên)
	RoleModerator                        // 'moderator' (Kiểm duyệt viên)
)