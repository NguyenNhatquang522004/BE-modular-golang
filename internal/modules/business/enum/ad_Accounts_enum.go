package enum

// =============================================================================
// AD ACCOUNT STATUS
// =============================================================================

//go:generate enumer -type=AccountStatus -json -transform=snake -trimprefix=AccountStatus
type AccountStatus int

const (
	AccountStatusActive   AccountStatus = iota // 'active'
	AccountStatusDisabled                      // 'disabled'
	AccountStatusSettled                       // 'settled' (Đã quyết toán công nợ)
)