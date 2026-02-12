package enum

//go:generate enumer -type=Type_Block -json -sql -transform=snake -trimprefix=Type_Block
type Type_Block int

const (
	Type_Block_Full Type_Block = iota
	Type_Block_Partial
)
