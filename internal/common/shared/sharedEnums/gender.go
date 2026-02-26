package sharedEnums
//go:generate enumer -type=Gender -json -transform=snake -trimprefix=Gender
type Gender int

const (
	GenderMale Gender = iota
	GenderFemale
	GenderOther
	GenderHidden
)

