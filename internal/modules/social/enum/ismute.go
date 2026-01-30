package enum

//go:generate enumer -type=Is_Muted -json -sql -transform=snake -trimprefix=Is_Muted
type Is_Muted int
const (
	Is_Muted_No Is_Muted = iota
	Is_Muted_Yes
)