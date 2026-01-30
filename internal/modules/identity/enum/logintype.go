package enum

// 1. Định nghĩa Type riêng (Clean & Type Safe)
//
//go:generate enumer -type=LoginType -json -sql -transform=snake -trimprefix=LoginType
type LoginType int

// 2. Định nghĩa các giá trị được phép (Single Source of Truth)
const (
	LoginTypeEmail LoginType = iota
	LoginTypeGoogle
	LoginTypeFacebook
	LoginTypeApple
)
