package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------
// MAPPER REQ TO ENTITY
// ---------------------------------------------------------

// ToEntityGroupEvent: Ánh xạ 100% từ CreateGroupEventReq sang Entity GroupEvent
func ToEntityGroupEvent(r *req.CreateGroupEventReq) *entity.GroupEvent {
	if r == nil {
		return nil
	}

	groupID, _ := primitive.ObjectIDFromHex(r.GroupID) // Cần validate GroupID ở tầng Transport/Handler trước

	return &entity.GroupEvent{
		ID:          primitive.NewObjectID(),
		GroupID:     groupID,
		CreatorID:   r.CreatorID,
		Title:       r.Title,
		Description: r.Description,
		CoverURL:    r.CoverURL,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		Location: entity.EventLocation{
			Type:        r.Location.Type,
			Address:     r.Location.Address,
			Coordinates: r.Location.Coordinates, // Copy mảng []float64
		},
		AttendeesCount: entity.EventAttendeeStats{
			Going:      r.AttendeesCount.Going,
			Interested: r.AttendeesCount.Interested,
		},
		CreatedAt: time.Now(),
	}
}

// UpdateToEntityGroupEvent: Cập nhật đè (Partial Update) từ UpdateGroupEventReq vào Entity hiện tại
func UpdateToEntityGroupEvent(r *req.UpdateGroupEventReq, e *entity.GroupEvent) {
	if r == nil || e == nil {
		return
	}

	if r.GroupID != nil {
		if oid, err := primitive.ObjectIDFromHex(*r.GroupID); err == nil {
			e.GroupID = oid
		}
	}
	if r.CreatorID != nil {
		e.CreatorID = *r.CreatorID
	}
	if r.Title != nil {
		e.Title = *r.Title
	}
	if r.Description != nil {
		e.Description = *r.Description
	}
	if r.CoverURL != nil {
		e.CoverURL = *r.CoverURL
	}
	if r.StartTime != nil {
		e.StartTime = *r.StartTime
	}
	if r.EndTime != nil {
		e.EndTime = *r.EndTime
	}
	if r.Location != nil {
		e.Location.Type = r.Location.Type
		e.Location.Address = r.Location.Address
		e.Location.Coordinates = r.Location.Coordinates
	}
	if r.AttendeesCount != nil {
		e.AttendeesCount.Going = r.AttendeesCount.Going
		e.AttendeesCount.Interested = r.AttendeesCount.Interested
	}
}

// ---------------------------------------------------------
// MAPPER ENTITY TO RES
// ---------------------------------------------------------

// ToResGroupEvent: Ánh xạ 100% từ Entity sang DTO Response (Chuyển ObjectID thành String)
func ToResGroupEvent(e *entity.GroupEvent) *res.GroupEventRes {
	if e == nil {
		return nil
	}

	return &res.GroupEventRes{
		ID:          e.ID.Hex(),
		GroupID:     e.GroupID.Hex(),
		CreatorID:   e.CreatorID,
		Title:       e.Title,
		Description: e.Description,
		CoverURL:    e.CoverURL,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location: res.ResEventLocation{
			Type:        e.Location.Type,
			Address:     e.Location.Address,
			Coordinates: e.Location.Coordinates,
		},
		AttendeesCount: res.ResEventAttendeeStats{
			Going:      e.AttendeesCount.Going,
			Interested: e.AttendeesCount.Interested,
		},
		CreatedAt: e.CreatedAt,
	}
}
