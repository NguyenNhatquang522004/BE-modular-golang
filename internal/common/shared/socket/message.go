package socket



// const (
// 	// System Events
// 	EventError     = "system_error"
// 	EventConnected = "system_connected"

// 	// Chat Events
// 	EventNewMessage = "chat_new_message"
// 	EventUserTyping = "chat_user_typing"
// 	EventReadMsg    = "chat_read_message"

// 	// Notification Events
// 	EventNewNotification = "notification_new"

// 	// Social Feed Events (Pub/Sub)
// 	EventNewComment = "feed_new_comment"
// 	EventNewLike    = "feed_new_like"

// 	// User Status Events
// 	EventUserOnline  = "user_online"
// 	EventUserOffline = "user_offline"
// )

// // ==========================================
// // 2. SERVER RESPONSE (Server -> Client)
// // ==========================================

// // Message là cấu trúc chuẩn được gửi từ Server xuống Client.
// // Sử dụng mô hình "Envelope" (Phong bì) để bọc dữ liệu.
// type Message struct {
// 	// Các trường dùng để định tuyến (Routing) - KHÔNG gửi xuống Client
// 	ToUserID string `json:"-"` // Chỉ gửi cho User cụ thể này
// 	Topic    string `json:"-"` // Chỉ gửi cho Room/Topic này (Pub/Sub)

// 	// Các trường dữ liệu thực tế gửi xuống Client
// 	Type      string      `json:"type"`                // Loại sự kiện (VD: "new_notification")
// 	Payload   interface{} `json:"payload"`             // Dữ liệu chi tiết (Struct hoặc Map)
// 	Timestamp int64       `json:"timestamp,omitempty"` // Thời gian gửi (Unix User)
// }

// // Hàm helper để tạo tin nhắn nhanh
// func NewMessage(eventType string, payload interface{}) Message {
// 	return Message{
// 		Type:      eventType,
// 		Payload:   payload,
// 		Timestamp: time.Now().UnixMilli(),
// 	}
// }

// // ==========================================
// // 3. CLIENT REQUEST (Client -> Server)
// // ==========================================

// // ClientRequest là cấu trúc tin nhắn mà Client gửi lên Server.
// // Trong mô hình Hybrid, Client chủ yếu dùng cái này để Subscribe/Unsubscribe Topic.
// type ClientRequest struct {
// 	Action string `json:"action"`         // VD: "subscribe", "unsubscribe", "typing"
// 	Topic  string `json:"topic"`          // VD: "post_123", "chat_room_456"
// 	Data   string `json:"data,omitempty"` // Dữ liệu phụ (nếu cần)
// }

// // Các Action chuẩn của Client
// const (
// 	ActionSubscribe   = "subscribe"
// 	ActionUnsubscribe = "unsubscribe"
// 	ActionPing        = "ping"
// )
