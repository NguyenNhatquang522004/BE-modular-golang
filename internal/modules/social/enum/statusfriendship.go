package enum

//go:generate enumer -type=StatusFriendship -json -sql -transform=snake -trimprefix=Type_Block

type StatusFriendship int

const (
	StatusFriendship_Pending StatusFriendship = iota
	StatusFriendship_Accepted
	StatusFriendship_Declined
	StatusFriendship_Blocked
)
