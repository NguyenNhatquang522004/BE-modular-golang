package enum

// --- 1. POST TYPE ---
//go:generate enumer -type=PostType -json -transform=snake -trimprefix=PostType
type PostType int

const (
	PostTypeText       PostType = iota // 'text'
	PostTypeMedia                      // 'media'
	PostTypeShare                      // 'share'
	PostTypeBackground                 // 'background'
	PostTypeQnA                        // 'qna'
	PostTypeLive                       // 'live'
)

// --- 2. CONTEXT TYPE ---
//go:generate enumer -type=ContextType -json -transform=snake -trimprefix=ContextType
type ContextType int

const (
	ContextTypeUserWall ContextType = iota // 'user_wall'
	ContextTypeGroup                       // 'group'
	ContextTypePage                        // 'page'
)

// --- 3. PRIVACY SCOPE ---
//go:generate enumer -type=PrivacyScope -json -transform=snake -trimprefix=Scope
type PrivacyScope int

const (
	ScopePublic      PrivacyScope = iota // 'public'
	ScopeFriends                         // 'friends'
	ScopeOnlyMe                          // 'only_me'
	ScopeCustom                          // 'custom'
)

// --- 4. STATUS ---
//go:generate enumer -type=PostStatus -json -transform=snake -trimprefix=Status
type PostStatus int

const (
	StatusPublished PostStatus = iota // 'published'
	StatusDraft                       // 'draft'
	StatusArchived                    // 'archived'
	StatusHidden                      // 'hidden'
)