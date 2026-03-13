package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// GroupEventRes - Phản hồi 100% dữ liệu từ Entity (ObjectID chuyển sang String)
type GroupEventRes struct {
	ID             string                `json:"id"`
	GroupID        string                `json:"group_id"`
	CreatorID      string                `json:"creator_id"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	CoverURL       string                `json:"cover_url"`
	StartTime      time.Time             `json:"start_time"`
	EndTime        time.Time             `json:"end_time"`
	Location       ResEventLocation      `json:"location"`
	AttendeesCount ResEventAttendeeStats `json:"attendees_count"`
	CreatedAt      time.Time             `json:"created_at"`
}

// --- Nested Structs ---

type ResEventLocation struct {
	Type        sharedEnums.EventLocationType `json:"type"`
	Address     string                        `json:"address"`
	Coordinates []float64                     `json:"coordinates,omitempty"`
}

type ResEventAttendeeStats struct {
	Going      int `json:"going"`
	Interested int `json:"interested"`
}
