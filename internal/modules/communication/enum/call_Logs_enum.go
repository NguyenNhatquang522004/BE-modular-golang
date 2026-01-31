package enum

// =============================================================================
// 1. CALL TYPE
// =============================================================================

//go:generate enumer -type=CallType -json -transform=snake -trimprefix=CallType
type CallType int

const (
	CallTypeVoice CallType = iota // 'voice'
	CallTypeVideo                 // 'video'
)

// =============================================================================
// 2. CALL STATUS
// =============================================================================

//go:generate enumer -type=CallStatus -json -transform=snake -trimprefix=CallStatus
type CallStatus int

const (
	CallStatusMissed   CallStatus = iota // 'missed' (Người nhận không nghe máy)
	CallStatusEnded                      // 'ended' (Cuộc gọi thành công và đã kết thúc)
	CallStatusRejected                   // 'rejected' (Người nhận bấm từ chối)
	CallStatusBusy                       // 'busy' (Người nhận đang bận cuộc gọi khác)
    // Có thể thêm: CallStatusFailed (Lỗi mạng), CallStatusCanceled (Người gọi dập máy trước khi đổ chuông)
)