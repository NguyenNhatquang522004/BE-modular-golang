package enum

//go:generate enumer -type=LiveStatus -json -transform=snake -trimprefix=LiveStatus
type LiveStatus int

const (
	LiveStatusPending   LiveStatus = iota // 'pending' (Đã lên lịch/Tạo phòng chờ)
	LiveStatusStreaming                   // 'streaming' (Đang phát)
	LiveStatusEnded                       // 'ended' (Đã kết thúc)
	LiveStatusBanned                      // 'banned' (Bị hệ thống chặn do vi phạm)
)