package enum

// =============================================================================
// SEARCH TARGET TYPE
// =============================================================================
// Loại đối tượng mà user đã click vào sau khi tìm kiếm.

//go:generate enumer -type=SearchTargetType -json -transform=snake -trimprefix=Target
type SearchTargetType int

const (
	TargetUser  SearchTargetType = iota // 'user'
	TargetGroup                         // 'group'
	TargetPage                          // 'page'
	TargetPost                          // 'post'
)