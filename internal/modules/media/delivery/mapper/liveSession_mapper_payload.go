package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreatePayloadToEntity(
	p *mediaEvent.CreateLiveSessionPayload,
	streamKey string,
	playbackURL string,
) *entity.LiveSession {
	now := time.Now().UTC()
	if p.SessionID != "" {
		return nil
	}
	// An toàn với con trỏ: Mặc dù đã có binding:"required",
	// best practice vẫn là check nil để chống panic 100% ở tầng logic.
	isRecorded := false
	if p.IsRecorded != nil {
		isRecorded = *p.IsRecorded
	}
	entity := &entity.LiveSession{
		ID:          primitive.NewObjectID(), // Tự sinh ObjectID cho Mongo
		HostUserID:  p.UserID,
		PageID:      *p.PageID,
		GroupID:     *p.GroupID,
		Title:       p.Title,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		Status:      sharedEnums.ProcessingPending, // Trạng thái mặc định

		// Technical Info
		StreamKey:   streamKey,
		PlaybackURL: playbackURL,
		RecordingSetting: entity.RecordingSetting{
			IsRecorded: isRecorded,
		},

		// Best Practice: Khởi tạo mảng rỗng để JSON trả ra [] thay vì null
		BannedUsers: make([]string, 0),

		CreatedAt: now,
		UpdatedAt: now,
	}
	return entity
}
func ApplyUpdatePayloadToEntity(existingLive *entity.LiveSession, p *mediaEvent.UpdateLiveSessionPayload) *entity.LiveSession {
	// Dùng cờ hasChanged để theo dõi xem có dữ liệu nào THỰC SỰ thay đổi không
	hasChanged := false

	// Kiểm tra Title
	if p.Title != nil && *p.Title != existingLive.Title {
		existingLive.Title = *p.Title
		hasChanged = true
	}

	// Kiểm tra Description
	if p.Description != nil && *p.Description != existingLive.Description {
		existingLive.Description = *p.Description
		hasChanged = true
	}

	// Kiểm tra CategoryID
	if p.CategoryID != nil && *p.CategoryID != existingLive.CategoryID {
		existingLive.CategoryID = *p.CategoryID
		hasChanged = true
	}

	// Kiểm tra RecordingSetting
	if p.IsRecorded != nil && *p.IsRecorded != existingLive.RecordingSetting.IsRecorded {
		existingLive.RecordingSetting.IsRecorded = *p.IsRecorded
		hasChanged = true
	}

	// CHỈ cập nhật UpdatedAt nếu có dữ liệu thực sự bị thay đổi
	if hasChanged {
		existingLive.UpdatedAt = time.Now().UTC()
	}

	// Trả về chính con trỏ entity đó sau khi đã update (giúp code ở Service chain dễ dàng hơn)
	return existingLive
}
