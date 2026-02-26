package sharedEnums

//go:generate enumer -type=PrivacyScope -json -transform=snake -trimprefix=Scope
type PrivacyScope int

const (
	ScopePublic  PrivacyScope = iota // 'public'
	ScopeFriends                     // 'friends'
	ScopeOnlyMe                      // 'only_me'
	ScopeCustom                      // 'custom'
)
