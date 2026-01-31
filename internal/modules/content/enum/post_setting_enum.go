package enum

//go:generate enumer -type=PublisherRole -json -transform=snake -trimprefix=Role
type PublisherRole int

const (
	RoleAdmin     PublisherRole = iota // 'admin'
	RoleEditor                         // 'editor'
	RoleModerator                      // 'moderator'
	RoleBot                            // 'bot' (Nếu là auto post)
)