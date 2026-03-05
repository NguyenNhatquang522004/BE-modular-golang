package sharedEnums

//go:generate enumer -type=EmailFrequency -json -transform=snake -trimprefix=Email
type EmailFrequency int

const (
	EmailDaily EmailFrequency = iota
	EmailWeekly
	EmailNever
)
