package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToReelEntity(payload *mediaEvent.CreateReelPayload, userID string) (*entity.Reel, error) {
	now := time.Now()

	// Khởi tạo Entity với các trường mặc định của hệ thống
	reel := &entity.Reel{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		ProcessingStatus: sharedEnums.ProcessingPending, // Giả sử enum của bạn có trạng thái Pending
		Caption:          payload.Caption,
		Hashtags:         payload.Hashtags,
		Mentions:         payload.Mentions,
		Privacy:          payload.Privacy,
		CreatedAt:        now,
		Stats:            entity.ReelStats{}, // Mặc định tất cả bằng 0
	}

	// 1. Map Video
	reel.Video = entity.ReelVideo{
		URL:           payload.Video.URL,
		ThumbnailURL:  payload.Video.ThumbnailURL,
		PreviewGifURL: payload.Video.PreviewGifURL,
		Width:         payload.Video.Width,
		Height:        payload.Video.Height,
		Duration:      payload.Video.Duration,
	}

	// 2. Map AudioMeta
	reel.AudioMeta = entity.AudioMeta{
		IsOriginalAudio: payload.AudioMeta.IsOriginalAudio,
		VolumeAdjust:    payload.AudioMeta.VolumeAdjust,
		AudioStartTime:  payload.AudioMeta.AudioStartTime,
	}

	// Kiểm tra nếu client có gửi TrackID thì tiến hành parse sang ObjectID
	if payload.AudioMeta.TrackID != nil && *payload.AudioMeta.TrackID != "" {
		trackID, err := primitive.ObjectIDFromHex(*payload.AudioMeta.TrackID)
		if err != nil {
			return nil, fmt.Errorf("invalid audio track_id format: %w", err)
		}
		reel.AudioMeta.TrackID = &trackID
	}

	// 3. Map RemixInfo (Nếu là video remix)
	if payload.RemixInfo != nil {
		parentID, err := primitive.ObjectIDFromHex(payload.RemixInfo.ParentReelID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_reel_id format: %w", err)
		}
		reel.RemixInfo = &entity.RemixInfo{
			ParentReelID: parentID,
			Type:         payload.RemixInfo.Type,
			IsRemixable:  payload.RemixInfo.IsRemixable,
		}
	}

	return reel, nil
}

func UpdateReelEntity(existingReel *entity.Reel, payload *mediaEvent.UpdateReelPayload) *entity.Reel {
	// Sử dụng Pointer từ DTO để phân biệt:
	// - nil: Client không muốn update trường này.
	// - != nil: Client muốn update (kể cả truyền string rỗng "").

	if payload.Caption != nil {
		existingReel.Caption = *payload.Caption
	}

	if payload.Privacy != nil {
		existingReel.Privacy = *payload.Privacy
	}

	// Với slice (mảng), JSON Unmarshal trong Go hoạt động như sau:
	// - Không gửi field -> nil
	// - Gửi mảng rỗng [] -> slice có len = 0, cap = 0 (nhưng != nil)
	// Do đó, ta chỉ cần check != nil là có thể update, cho phép client xóa hết hashtag bằng cách gửi []
	if payload.Hashtags != nil {
		existingReel.Hashtags = payload.Hashtags
	}

	if payload.Mentions != nil {
		existingReel.Mentions = payload.Mentions
	}

	// Lưu ý: Thường Entity sẽ có thêm trường UpdatedAt.
	// Nếu Entity của bạn thêm trường này trong tương lai, bạn sẽ gán:
	// existingReel.UpdatedAt = time.Now() ở ngay đây.

	return existingReel
}
