package enum

// =============================================================================
// 1. CAMPAIGN OBJECTIVE (Mục tiêu chiến dịch)
// =============================================================================

//go:generate enumer -type=CampaignObjective -json -transform=snake -trimprefix=Objective
type CampaignObjective int

const (
	ObjectiveReach       CampaignObjective = iota // 'reach'
	ObjectiveTraffic                              // 'traffic'
	ObjectiveMessages                             // 'messages'
	ObjectiveConversions                          // 'conversions'
)

// =============================================================================
// 2. BUYING TYPE (Hình thức đấu thầu)
// =============================================================================

//go:generate enumer -type=BuyingType -json -transform=snake -trimprefix=Buying
type BuyingType int

const (
	BuyingAuction    BuyingType = iota // 'auction' (Đấu giá)
	BuyingFixedPrice                   // 'fixed_price' (Giá cố định/Reservation)
)

// =============================================================================
// 3. CAMPAIGN STATUS
// =============================================================================

//go:generate enumer -type=CampaignStatus -json -transform=snake -trimprefix=CampaignStatus
type CampaignStatus int

const (
	CampaignStatusActive    CampaignStatus = iota // 'active'
	CampaignStatusPaused                          // 'paused'
	CampaignStatusCompleted                       // 'completed'
	CampaignStatusArchived                        // 'archived'
)