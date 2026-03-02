package enum

// =============================================================================
// 1. PAGE STATUS
// =============================================================================

//go:generate enumer -type=PageStatus -json -transform=snake -trimprefix=PageStatus
type PageStatus int

const (
	PageStatusPublished   PageStatus = iota // 'published'
	PageStatusUnpublished                   // 'unpublished' (Chủ page ẩn)
	PageStatusBanned                        // 'banned' (Admin hệ thống chặn)
)

// =============================================================================
// 2. CTA TYPE (Call To Action)
// =============================================================================

//go:generate enumer -type=CTAType -json -transform=snake -trimprefix=CTA
type CTAType int

const (
	CTASendMessage CTAType = iota // 'send_message'
	CTACallNow                    // 'call_now'
	CTAVisitWebsite               // 'visit_website'
	CTAShopNow                    // 'shop_now'
	CTAFollow                     // 'follow'
)


//go:generate enumer -type=MessagingStatus -json -transform=snake -trimprefix=MsgStatus
type MessagingStatus int

const (
	MsgStatusOnline MessagingStatus = iota // 'online'
	MsgStatusAway                          // 'away'
)

// =============================================================================
// 4. DAY OF WEEK (Cho Business Hours)
// =============================================================================

//go:generate enumer -type=DayOfWeek -json -transform=title -trimprefix=Day
type DayOfWeek int

const (
	DayMonday DayOfWeek = iota
	DayTuesday
	DayWednesday
	DayThursday
	DayFriday
	DaySaturday
	DaySunday
)