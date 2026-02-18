package enum

//go:generate enumer -type=Type_Block -json -sql -transform=snake -trimprefix=Type_Block
type Type_Block int

const (
	Type_Block_Full Type_Block = iota
	Type_Block_Partial
	Type_Block_Chat
	Type_Block_Profile
	Type_Block_None
	//... Thêm các loại block khác vào đây
)
