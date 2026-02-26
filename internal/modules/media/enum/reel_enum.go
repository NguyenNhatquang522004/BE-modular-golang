package enum

// =============================================================================
// 2. REMIX TYPE
// =============================================================================

//go:generate enumer -type=RemixType -json -transform=snake -trimprefix=Remix
type RemixType int

const (
	RemixDuet  RemixType = iota // 'duet' (Chia đôi màn hình)
	RemixRemix                  // 'remix' (Lồng ghép)
)



