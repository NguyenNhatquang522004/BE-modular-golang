package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// CreateGroupEventReq - Ánh xạ 100% các trường đầu vào cần thiết để tạo sự kiện
type CreateGroupEventReq struct {
	GroupID        string                `json:"group_id" binding:"required"`
	CreatorID      string                `json:"creator_id" binding:"required"`
	Title          string                `json:"title" binding:"required"`
	Description    string                `json:"description"`
	CoverURL       string                `json:"cover_url"`
	StartTime      time.Time             `json:"start_time" binding:"required"`
	EndTime        time.Time             `json:"end_time" binding:"required,gtfield=StartTime"` // EndTime phải lớn hơn StartTime
	Location       ReqEventLocation      `json:"location" binding:"required"`
	AttendeesCount ReqEventAttendeeStats `json:"attendees_count"` // Thông thường mặc định là 0 khi tạo, nhưng map đủ 100% theo yêu cầu
}

// UpdateGroupEventReq - Dùng con trỏ 100% để hỗ trợ Partial Update
type UpdateGroupEventReq struct {
	GroupID        *string                `json:"group_id,omitempty"`
	CreatorID      *string                `json:"creator_id,omitempty"`
	Title          *string                `json:"title,omitempty"`
	Description    *string                `json:"description,omitempty"`
	CoverURL       *string                `json:"cover_url,omitempty"`
	StartTime      *time.Time             `json:"start_time,omitempty"`
	EndTime        *time.Time             `json:"end_time,omitempty"`
	Location       *ReqEventLocation      `json:"location,omitempty"`
	AttendeesCount *ReqEventAttendeeStats `json:"attendees_count,omitempty"`
}

// --- Nested Structs ---

type ReqEventLocation struct {
	Type        sharedEnums.EventLocationType `json:"type" binding:"required"`
	Address     string                        `json:"address"`
	Coordinates []float64                     `json:"coordinates,omitempty"` // [Kinh độ, Vĩ độ]
}

type ReqEventAttendeeStats struct {
	Going      int `json:"going"`
	Interested int `json:"interested"`
}
