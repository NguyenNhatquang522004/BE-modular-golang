package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntityReel: Chuyển từ Req -> Entity để Insert
func ToEntityReel(r *req.ReelReq) (*entity.Reel, error) {
	if r == nil {
		return nil, nil
	}

	e := &entity.Reel{
		UserID:           r.UserID,
		ProcessingStatus: r.ProcessingStatus,
		Caption:          r.Caption,
		Hashtags:         r.Hashtags,
		Mentions:         r.Mentions,
		Privacy:          r.Privacy,
		Video: entity.ReelVideo{
			URL:           r.Video.URL,
			ThumbnailURL:  r.Video.ThumbnailURL,
			PreviewGifURL: r.Video.PreviewGifURL,
			Width:         r.Video.Width,
			Height:        r.Video.Height,
			Duration:      r.Video.Duration,
		},
		AudioMeta: entity.AudioMeta{
			IsOriginalAudio: r.AudioMeta.IsOriginalAudio,
			VolumeAdjust:    r.AudioMeta.VolumeAdjust,
			AudioStartTime:  r.AudioMeta.AudioStartTime,
		},
		Stats: entity.ReelStats{
			Views:    r.Stats.Views,
			Likes:    r.Stats.Likes,
			Shares:   r.Stats.Shares,
			Saves:    r.Stats.Saves,
			Comments: r.Stats.Comments,
		},
		DeletedAt: r.DeletedAt,
	}

	// Xử lý ID chính
	if r.ID != "" {
		id, err := primitive.ObjectIDFromHex(r.ID)
		if err == nil {
			e.ID = id
		}
	} else {
		e.ID = primitive.NewObjectID()
	}

	// Xử lý AudioMeta TrackID (Pointer)
	if r.AudioMeta.TrackID != "" {
		if trackID, err := primitive.ObjectIDFromHex(r.AudioMeta.TrackID); err == nil {
			e.AudioMeta.TrackID = &trackID
		}
	}

	// Xử lý RemixInfo (Pointer)
	if r.RemixInfo != nil {
		remix := &entity.RemixInfo{
			Type:        r.RemixInfo.Type,
			IsRemixable: r.RemixInfo.IsRemixable,
		}
		if r.RemixInfo.ParentReelID != "" {
			if pID, err := primitive.ObjectIDFromHex(r.RemixInfo.ParentReelID); err == nil {
				remix.ParentReelID = pID
			}
		}
		e.RemixInfo = remix
	}

	// Xử lý Timestamps
	now := time.Now()
	if r.CreatedAt != nil {
		e.CreatedAt = *r.CreatedAt
	} else {
		e.CreatedAt = now
	}

	return e, nil
}

// 2. UpdateToEntityReel: Cập nhật Partial Update
func UpdateToEntityReel(r *req.UpdateReelReq, e *entity.Reel) {
	if r == nil || e == nil {
		return
	}

	if r.ProcessingStatus != nil {
		e.ProcessingStatus = *r.ProcessingStatus
	}
	if r.Caption != nil {
		e.Caption = *r.Caption
	}
	if r.Hashtags != nil {
		e.Hashtags = r.Hashtags
	}
	if r.Mentions != nil {
		e.Mentions = r.Mentions
	}
	if r.Privacy != nil {
		e.Privacy = *r.Privacy
	}
	if r.DeletedAt != nil {
		e.DeletedAt = r.DeletedAt
	}

	// Update Video
	if r.Video != nil {
		e.Video.URL = r.Video.URL
		e.Video.ThumbnailURL = r.Video.ThumbnailURL
		e.Video.PreviewGifURL = r.Video.PreviewGifURL
		e.Video.Width = r.Video.Width
		e.Video.Height = r.Video.Height
		e.Video.Duration = r.Video.Duration
	}

	// Update AudioMeta
	if r.AudioMeta != nil {
		e.AudioMeta.IsOriginalAudio = r.AudioMeta.IsOriginalAudio
		e.AudioMeta.VolumeAdjust = r.AudioMeta.VolumeAdjust
		e.AudioMeta.AudioStartTime = r.AudioMeta.AudioStartTime
		if r.AudioMeta.TrackID != "" {
			if trackID, err := primitive.ObjectIDFromHex(r.AudioMeta.TrackID); err == nil {
				e.AudioMeta.TrackID = &trackID
			}
		}
	}

	// Update RemixInfo
	if r.RemixInfo != nil {
		if e.RemixInfo == nil {
			e.RemixInfo = &entity.RemixInfo{}
		}
		e.RemixInfo.Type = r.RemixInfo.Type
		e.RemixInfo.IsRemixable = r.RemixInfo.IsRemixable
		if r.RemixInfo.ParentReelID != "" {
			if pID, err := primitive.ObjectIDFromHex(r.RemixInfo.ParentReelID); err == nil {
				e.RemixInfo.ParentReelID = pID
			}
		}
	}

	// Update Stats
	if r.Stats != nil {
		e.Stats.Views = r.Stats.Views
		e.Stats.Likes = r.Stats.Likes
		e.Stats.Shares = r.Stats.Shares
		e.Stats.Saves = r.Stats.Saves
		e.Stats.Comments = r.Stats.Comments
	}
}

// 3. ReqToResReel: Chuyển thẳng từ Req -> Res
func ReqToResReel(r *req.ReelReq) *res.ReelRes {
	if r == nil {
		return nil
	}

	response := &res.ReelRes{
		ID:               r.ID,
		UserID:           r.UserID,
		ProcessingStatus: r.ProcessingStatus,
		Caption:          r.Caption,
		Hashtags:         r.Hashtags,
		Mentions:         r.Mentions,
		Privacy:          r.Privacy,
		Video: res.ReelVideoRes{
			URL:           r.Video.URL,
			ThumbnailURL:  r.Video.ThumbnailURL,
			PreviewGifURL: r.Video.PreviewGifURL,
			Width:         r.Video.Width,
			Height:        r.Video.Height,
			Duration:      r.Video.Duration,
		},
		AudioMeta: res.AudioMetaRes{
			TrackID:         r.AudioMeta.TrackID,
			IsOriginalAudio: r.AudioMeta.IsOriginalAudio,
			VolumeAdjust:    r.AudioMeta.VolumeAdjust,
			AudioStartTime:  r.AudioMeta.AudioStartTime,
		},
		Stats: res.ReelStatsRes{
			Views:    r.Stats.Views,
			Likes:    r.Stats.Likes,
			Shares:   r.Stats.Shares,
			Saves:    r.Stats.Saves,
			Comments: r.Stats.Comments,
		},
		DeletedAt: r.DeletedAt,
	}

	if r.RemixInfo != nil {
		response.RemixInfo = &res.RemixInfoRes{
			ParentReelID: r.RemixInfo.ParentReelID,
			Type:         r.RemixInfo.Type,
			IsRemixable:  r.RemixInfo.IsRemixable,
		}
	}

	if r.CreatedAt != nil {
		response.CreatedAt = *r.CreatedAt
	} else {
		response.CreatedAt = time.Now()
	}

	return response
}

// 4. EntityToResReel: Chuyển dữ liệu từ DB -> Chuẩn DTO Response cho Client
func EntityToResReel(e *entity.Reel) *res.ReelRes {
	if e == nil {
		return nil
	}

	response := &res.ReelRes{
		ID:               e.ID.Hex(),
		UserID:           e.UserID,
		ProcessingStatus: e.ProcessingStatus,
		Caption:          e.Caption,
		Hashtags:         e.Hashtags,
		Mentions:         e.Mentions,
		Privacy:          e.Privacy,
		Video: res.ReelVideoRes{
			URL:           e.Video.URL,
			ThumbnailURL:  e.Video.ThumbnailURL,
			PreviewGifURL: e.Video.PreviewGifURL,
			Width:         e.Video.Width,
			Height:        e.Video.Height,
			Duration:      e.Video.Duration,
		},
		AudioMeta: res.AudioMetaRes{
			IsOriginalAudio: e.AudioMeta.IsOriginalAudio,
			VolumeAdjust:    e.AudioMeta.VolumeAdjust,
			AudioStartTime:  e.AudioMeta.AudioStartTime,
		},
		Stats: res.ReelStatsRes{
			Views:    e.Stats.Views,
			Likes:    e.Stats.Likes,
			Shares:   e.Stats.Shares,
			Saves:    e.Stats.Saves,
			Comments: e.Stats.Comments,
		},
		CreatedAt: e.CreatedAt,
		DeletedAt: e.DeletedAt,
	}

	if e.AudioMeta.TrackID != nil && !e.AudioMeta.TrackID.IsZero() {
		response.AudioMeta.TrackID = e.AudioMeta.TrackID.Hex()
	}

	if e.RemixInfo != nil {
		response.RemixInfo = &res.RemixInfoRes{
			Type:        e.RemixInfo.Type,
			IsRemixable: e.RemixInfo.IsRemixable,
		}
		if !e.RemixInfo.ParentReelID.IsZero() {
			response.RemixInfo.ParentReelID = e.RemixInfo.ParentReelID.Hex()
		}
	}

	return response
}
