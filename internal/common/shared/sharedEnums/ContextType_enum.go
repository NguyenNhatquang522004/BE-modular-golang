package sharedEnums

//go:generate enumer -type=ContextType -json -transform=snake -trimprefix=ContextType
type ContextType int

const (
	ContextTypeUserWall ContextType = iota // 'user_wall'
	ContextTypeGroup                       // 'group'
	ContextTypePage                        // 'page'
)
