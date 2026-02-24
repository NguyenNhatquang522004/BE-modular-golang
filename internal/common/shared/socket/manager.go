package socket

import "context"

//	Interface để các module khác gọi
//
// Event types (Ví dụ: tin nhắn mới, thông báo mới)
const (
	EventNewMessage      = "new_message"
	EventNewNotification = "new_notification"
)

// Payload chuẩn cho mọi message gửi qua socket
type Message struct {
	ToUserID string      `json:"-"` // Chỉ dùng nội bộ để định tuyến
	Type     string      `json:"type"`
	Payload  interface{} `json:"payload"`
}

// Interface để các module khác inject vào
type Manager interface {
	Run(ctx context.Context) // Thêm Context để Graceful Shutdown
	Register(client ClientOps)
	Unregister(client ClientOps)
	// Hàm quan trọng nhất để các module khác gọi
	SendToUser(userID string, msg Message) error
	SendBulkToDirect(userIDs []string, msg Message) error
	Broadcast(msg Message)
}

// Interface để trừu tượng hóa Client
type ClientOps interface {
	ReadPump()
	WritePump()
	GetUserID() string
	Send(msg []byte)
	Close()
}
