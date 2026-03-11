package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToStoryEntity chuyển đổi payload tạo mới thành Story Entity hoàn chỉnh
func ToStoryEntity(payload *mediaEvent.CreateStoryPayload) *entity.Story {
	now := time.Now().UTC()

	// Tính toán thời gian hết hạn (24 giờ kể từ lúc đăng)
	expiresAt := now.Add(24 * time.Hour)

	// Xử lý Overlays (Khởi tạo capacity để tối ưu bộ nhớ)
	var overlays []*entity.StoryOverlay
	if len(payload.Overlays) > 0 {
		overlays = make([]*entity.StoryOverlay, 0, len(payload.Overlays))
		for _, o := range payload.Overlays {
			overlays = append(overlays, &entity.StoryOverlay{
				Type: o.Type,
				Position: entity.OverlayPosition{
					X:        o.Position.X,
					Y:        o.Position.Y,
					Rotation: o.Position.Rotation,
					Scale:    o.Position.Scale,
				},
				Data: o.Data,
			})
		}
	}

	// Đảm bảo AllowList và BlockList không bị nil (nếu rỗng thì gán mảng rỗng)
	// Tránh lưu giá trị null vào DB
	allowList := payload.Privacy.AllowList
	if allowList == nil {
		allowList = []string{}
	}

	blockList := payload.Privacy.BlockList
	if blockList == nil {
		blockList = []string{}
	}

	return &entity.Story{
		ID:     primitive.NewObjectID(),
		UserID: payload.UserID,
		Media: entity.StoryMedia{
			URL:          payload.Media.URL,
			Type:         payload.Media.Type,
			Duration:     payload.Media.Duration,
			ThumbnailURL: payload.Media.ThumbnailURL,
			SizeBytes:    payload.Media.SizeBytes,
			Width:        payload.Media.Width,
			Height:       payload.Media.Height,
			MimeType:     payload.Media.MimeType,
		},
		Overlays: overlays,
		Privacy: entity.StoryPrivacy{
			Type:      payload.Privacy.Type,
			AllowList: allowList,
			BlockList: blockList,
		},
		Settings: entity.StorySettings{
			AllowReply: payload.Settings.AllowReply,
			AllowShare: payload.Settings.AllowShare,
		},

		// Khởi tạo các giá trị mặc định cho Story mới
		PreviewViewers: make([]entity.ViewerPreview, 0), // Mảng rỗng thay vì nil
		Stats: entity.StoryStats{
			ViewsCount: 0,
			Likes:      0,
			Love:       0,
			Haha:       0,
			Wow:        0,
			Sad:        0,
			Angry:      0,
			ReplyCount: 0,
		},

		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		IsArchived: false,
	}
}

// =========================================================================
// 2. UPDATE MAPPER: DTO -> BSON Update Document
// =========================================================================

// BuildUpdateStoryBSON tạo ra cấu trúc map cho lệnh `$set` của MongoDB.
// Cách tiếp cận này ngăn chặn tình trạng Race Condition (ghi đè stats khi update setting)
func ApplyUpdateStoryPayload(story *entity.Story, payload *mediaEvent.UpdateStoryPayload) *entity.Story {
	// Kiểm tra an toàn
	if story == nil {
		return nil
	}

	// 1. Map Update Privacy
	if payload.Privacy != nil {
		if payload.Privacy.Type != nil {
			story.Privacy.Type = *payload.Privacy.Type
		}
		// Lưu ý: Nếu Client gửi mảng rỗng []string{}, nó vẫn được tính là khác nil
		// và sẽ ghi đè mảng cũ thành mảng rỗng (đúng với logic xóa danh sách)
		if payload.Privacy.AllowList != nil {
			story.Privacy.AllowList = payload.Privacy.AllowList
		}
		if payload.Privacy.BlockList != nil {
			story.Privacy.BlockList = payload.Privacy.BlockList
		}
	}

	// 2. Map Update Settings
	if payload.Settings != nil {
		if payload.Settings.AllowReply != nil {
			story.Settings.AllowReply = *payload.Settings.AllowReply
		}
		if payload.Settings.AllowShare != nil {
			story.Settings.AllowShare = *payload.Settings.AllowShare
		}
	}

	return story
}
