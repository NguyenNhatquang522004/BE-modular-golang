package sharedEnums

//go:generate enumer -type=RoleType -json -transform=snake -trimprefix=RoleType
type RoleType int

// 2. Định nghĩa các giá trị được phép (Single Source of Truth)
const (
	RoleTypeUser RoleType = iota
	RoleTypeAdmin
	RoleTypeSystemAdmin
	RoleTypeModerator
	RoleTypeMember
	RoleTypeGuest
	RoleTypeAdvertiser
	RoleTypeAnalyst
)
