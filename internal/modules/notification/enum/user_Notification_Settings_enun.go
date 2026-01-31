package enum

// =============================================================================
// EMAIL FREQUENCY
// =============================================================================

//go:generate enumer -type=EmailFrequency -json -transform=snake -trimprefix=EmailFreq
type EmailFrequency int

const (
	EmailFreqInstant EmailFrequency = iota // 'instant' (Gửi ngay lập tức - Mặc định)
	EmailFreqDaily                         // 'daily' (Tổng hợp 1 ngày 1 email)
	EmailFreqWeekly                        // 'weekly' (Tổng hợp 1 tuần 1 email)
	EmailFreqNever                         // 'never' (Không nhận email marketing/digest)
)