package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
	CollectionGroupEvents = "GroupEvents"
)
type GroupEvent struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LINKING
	// Sự kiện thuộc về nhóm nào (Mongo ID)
	// Index: { group_id: 1, start_time: 1 } -> Lấy lịch sự kiện của nhóm
	GroupID primitive.ObjectID `bson:"group_id" json:"group_id"`

	// Người tạo sự kiện (Postgres UUID -> String)
	CreatorID string `bson:"creator_id" json:"creator_id"`

	// 2. CONTENT
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	CoverURL    string `bson:"cover_url" json:"cover_url"`

	// 3. TIME
	// Index: { start_time: 1 } -> Lọc sự kiện sắp diễn ra
	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`

	// 4. LOCATION
	// Chứa thông tin địa điểm và tọa độ
	Location EventLocation `bson:"location" json:"location"`

	// 5. STATS
	AttendeesCount EventAttendeeStats `bson:"attendees_count" json:"attendees_count"`

	// 6. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// --- LOCATION ---
type EventLocation struct {
	Type enum.EventLocationType `bson:"type" json:"type"` // 'online', 'offline'

	// Địa chỉ text hoặc Link Online (Zoom/Google Meet)
	Address string `bson:"address" json:"address"`

	// GeoJSON Coordinates: [Longitude, Latitude] -> [Kinh độ, Vĩ độ]
	// Dùng []float64 để mapping chuẩn với mảng số của Mongo.
	// Index: 2dsphere
	Coordinates []float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

// --- ATTENDEES STATS ---
// Denormalization: Lưu số lượng để hiển thị nhanh.
// Danh sách chi tiết người tham gia sẽ nằm ở collection riêng (Event_Attendees)
type EventAttendeeStats struct {
	Going      int `bson:"going" json:"going"`           // Số người xác nhận đi
	Interested int `bson:"interested" json:"interested"` // Số người quan tâm
}
func (GroupEvent) CollectionName() string {
	return CollectionGroupEvents
}