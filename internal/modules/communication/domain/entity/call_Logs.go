package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
	CollectionCallLogs = "CallLogs"
)
type CallLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. CONTEXT
	// ID nhóm chat/hội thoại (Mongo ID)
	// Index: { conversation_id: 1, started_at: -1 } -> Lấy lịch sử call trong 1 nhóm
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`

	// 2. PARTICIPANTS
	// Người khởi tạo cuộc gọi (Postgres UUID -> String)
	CallerID string `bson:"caller_id" json:"caller_id"`

	// Danh sách người tham gia (bao gồm cả người gọi và người nhận)
	// Postgres UUID -> String
	// Index: Multikey { participants: 1, started_at: -1 } -> Lấy lịch sử call của User bất kỳ
	Participants []string `bson:"participants" json:"participants"`

	// 3. META INFO
	Type   enum.CallType   `bson:"type" json:"type"`     // 'voice', 'video'
	Status enum.CallStatus `bson:"status" json:"status"` // 'missed', 'ended'...

	// 4. TIMING
	StartedAt time.Time `bson:"started_at" json:"started_at"`

	// Có thể null nếu cuộc gọi không kết nối được hoặc logic xử lý lỗi
	EndedAt *time.Time `bson:"ended_at,omitempty" json:"ended_at,omitempty"`

	// Thời lượng tính bằng giây (0 nếu missed/rejected)
	DurationSeconds int `bson:"duration_seconds" json:"duration_seconds"`

	// 5. GROUP CALL FLAG
	IsGroupCall bool `bson:"is_group_call" json:"is_group_call"`
}
func (CallLog) CollectionName() string {
	return CollectionCallLogs
}