package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntity: Chuyển từ Req -> Entity để Insert
func ToEntityStory(r *req.StoryReq) (*entity.Story, error) {
	if r == nil {
		return nil, nil
	}

	e := &entity.Story{
		UserID: r.UserID,
		Media: entity.StoryMedia{
			URL:          r.Media.URL,
			Type:         r.Media.Type,
			Duration:     r.Media.Duration,
			ThumbnailURL: r.Media.ThumbnailURL,
			SizeBytes:    r.Media.SizeBytes,
		},
		Privacy: entity.StoryPrivacy{
			Type:      r.Privacy.Type,
			AllowList: r.Privacy.AllowList,
			BlockList: r.Privacy.BlockList,
		},
		Settings: entity.StorySettings{
			AllowReply: r.Settings.AllowReply,
			AllowShare: r.Settings.AllowShare,
		},
		Stats: entity.StoryStats{
			ViewsCount: r.Stats.ViewsCount,
			Likes:      r.Stats.Likes,
			ReplyCount: r.Stats.ReplyCount,
			Love:       r.Stats.Love,
			Haha:       r.Stats.Haha,
			Wow:        r.Stats.Wow,
			Sad:        r.Stats.Sad,
			Angry:      r.Stats.Angry,
		},
		IsArchived: r.IsArchived,
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

	// Xử lý mảng Overlays (Pointer Slice)
	if len(r.Overlays) > 0 {
		e.Overlays = make([]*entity.StoryOverlay, len(r.Overlays))
		for i, o := range r.Overlays {
			if o != nil {
				e.Overlays[i] = &entity.StoryOverlay{
					Type: o.Type,
					Position: entity.OverlayPosition{
						X:        o.Position.X,
						Y:        o.Position.Y,
						Rotation: o.Position.Rotation,
						Scale:    o.Position.Scale,
					},
					Data: o.Data, // Map copy reference
				}
			}
		}
	}

	// Xử lý mảng PreviewViewers
	if len(r.PreviewViewers) > 0 {
		e.PreviewViewers = make([]entity.ViewerPreview, len(r.PreviewViewers))
		for i, v := range r.PreviewViewers {
			e.PreviewViewers[i] = entity.ViewerPreview{
				UserID: v.UserID,
				Avatar: v.Avatar,
				Name:   v.Name,
			}
		}
	}

	// Xử lý Timestamps
	now := time.Now()
	if r.CreatedAt != nil {
		e.CreatedAt = *r.CreatedAt
	} else {
		e.CreatedAt = now
	}

	// Mặc định hết hạn sau 24h nếu không được chỉ định
	if r.ExpiresAt != nil {
		e.ExpiresAt = *r.ExpiresAt
	} else {
		e.ExpiresAt = e.CreatedAt.Add(24 * time.Hour)
	}

	return e, nil
}

// 2. UpdateToEntity: Cập nhật Partial Update
func UpdateToEntityStory(r *req.UpdateStoryReq, e *entity.Story) {
	if r == nil || e == nil {
		return
	}

	if r.Media != nil {
		e.Media.URL = r.Media.URL
		e.Media.Type = r.Media.Type
		e.Media.Duration = r.Media.Duration
		e.Media.ThumbnailURL = r.Media.ThumbnailURL
		e.Media.SizeBytes = r.Media.SizeBytes
	}

	if r.Overlays != nil { // Ghi đè toàn bộ danh sách overlays
		e.Overlays = make([]*entity.StoryOverlay, len(r.Overlays))
		for i, o := range r.Overlays {
			if o != nil {
				e.Overlays[i] = &entity.StoryOverlay{
					Type: o.Type,
					Position: entity.OverlayPosition{
						X:        o.Position.X,
						Y:        o.Position.Y,
						Rotation: o.Position.Rotation,
						Scale:    o.Position.Scale,
					},
					Data: o.Data,
				}
			}
		}
	}

	if r.Privacy != nil {
		e.Privacy.Type = r.Privacy.Type
		e.Privacy.AllowList = r.Privacy.AllowList
		e.Privacy.BlockList = r.Privacy.BlockList
	}

	if r.Settings != nil {
		e.Settings.AllowReply = r.Settings.AllowReply
		e.Settings.AllowShare = r.Settings.AllowShare
	}

	if r.PreviewViewers != nil { // Ghi đè mảng preview cache
		e.PreviewViewers = make([]entity.ViewerPreview, len(r.PreviewViewers))
		for i, v := range r.PreviewViewers {
			e.PreviewViewers[i] = entity.ViewerPreview{
				UserID: v.UserID,
				Avatar: v.Avatar,
				Name:   v.Name,
			}
		}
	}

	if r.Stats != nil {
		e.Stats.ViewsCount = r.Stats.ViewsCount
		e.Stats.Likes = r.Stats.Likes
		e.Stats.ReplyCount = r.Stats.ReplyCount
		e.Stats.Love = r.Stats.Love
		e.Stats.Haha = r.Stats.Haha
		e.Stats.Wow = r.Stats.Wow
		e.Stats.Sad = r.Stats.Sad
		e.Stats.Angry = r.Stats.Angry
	}

	if r.ExpiresAt != nil {
		e.ExpiresAt = *r.ExpiresAt
	}
	if r.IsArchived != nil {
		e.IsArchived = *r.IsArchived
	}
}

// 3. ReqToRes: Chuyển thẳng từ Req -> Res
func ReqToResStory(r *req.StoryReq) *res.StoryRes {
	if r == nil {
		return nil
	}

	response := &res.StoryRes{
		ID:     r.ID,
		UserID: r.UserID,
		Media: res.StoryMediaRes{
			URL:          r.Media.URL,
			Type:         r.Media.Type,
			Duration:     r.Media.Duration,
			ThumbnailURL: r.Media.ThumbnailURL,
			SizeBytes:    r.Media.SizeBytes,
		},
		Privacy: res.StoryPrivacyRes{
			Type:      r.Privacy.Type,
			AllowList: r.Privacy.AllowList,
			BlockList: r.Privacy.BlockList,
		},
		Settings: res.StorySettingsRes{
			AllowReply: r.Settings.AllowReply,
			AllowShare: r.Settings.AllowShare,
		},
		Stats: res.StoryStatsRes{
			ViewsCount: r.Stats.ViewsCount,
			Likes:      r.Stats.Likes,
			ReplyCount: r.Stats.ReplyCount,
			Love:       r.Stats.Love,
			Haha:       r.Stats.Haha,
			Wow:        r.Stats.Wow,
			Sad:        r.Stats.Sad,
			Angry:      r.Stats.Angry,
		},
		IsArchived: r.IsArchived,
	}

	if len(r.Overlays) > 0 {
		response.Overlays = make([]*res.StoryOverlayRes, len(r.Overlays))
		for i, o := range r.Overlays {
			if o != nil {
				response.Overlays[i] = &res.StoryOverlayRes{
					Type: o.Type,
					Position: res.OverlayPositionRes{
						X:        o.Position.X,
						Y:        o.Position.Y,
						Rotation: o.Position.Rotation,
						Scale:    o.Position.Scale,
					},
					Data: o.Data,
				}
			}
		}
	}

	if len(r.PreviewViewers) > 0 {
		response.PreviewViewers = make([]res.ViewerPreviewRes, len(r.PreviewViewers))
		for i, v := range r.PreviewViewers {
			response.PreviewViewers[i] = res.ViewerPreviewRes{
				UserID: v.UserID,
				Avatar: v.Avatar,
				Name:   v.Name,
			}
		}
	}

	now := time.Now()
	if r.CreatedAt != nil {
		response.CreatedAt = *r.CreatedAt
	} else {
		response.CreatedAt = now
	}

	if r.ExpiresAt != nil {
		response.ExpiresAt = *r.ExpiresAt
	} else {
		response.ExpiresAt = response.CreatedAt.Add(24 * time.Hour)
	}

	return response
}

// 4. EntityToRes: Sử dụng khi query DB -> Response Client
func EntityToResStory(e *entity.Story) *res.StoryRes {
	if e == nil {
		return nil
	}

	response := &res.StoryRes{
		ID:     e.ID.Hex(),
		UserID: e.UserID,
		Media: res.StoryMediaRes{
			URL:          e.Media.URL,
			Type:         e.Media.Type,
			Duration:     e.Media.Duration,
			ThumbnailURL: e.Media.ThumbnailURL,
			SizeBytes:    e.Media.SizeBytes,
		},
		Privacy: res.StoryPrivacyRes{
			Type:      e.Privacy.Type,
			AllowList: e.Privacy.AllowList,
			BlockList: e.Privacy.BlockList,
		},
		Settings: res.StorySettingsRes{
			AllowReply: e.Settings.AllowReply,
			AllowShare: e.Settings.AllowShare,
		},
		Stats: res.StoryStatsRes{
			ViewsCount: e.Stats.ViewsCount,
			Likes:      e.Stats.Likes,
			ReplyCount: e.Stats.ReplyCount,
			Love:       e.Stats.Love,
			Haha:       e.Stats.Haha,
			Wow:        e.Stats.Wow,
			Sad:        e.Stats.Sad,
			Angry:      e.Stats.Angry,
		},
		CreatedAt:  e.CreatedAt,
		ExpiresAt:  e.ExpiresAt,
		IsArchived: e.IsArchived,
	}

	if len(e.Overlays) > 0 {
		response.Overlays = make([]*res.StoryOverlayRes, len(e.Overlays))
		for i, o := range e.Overlays {
			if o != nil {
				response.Overlays[i] = &res.StoryOverlayRes{
					Type: o.Type,
					Position: res.OverlayPositionRes{
						X:        o.Position.X,
						Y:        o.Position.Y,
						Rotation: o.Position.Rotation,
						Scale:    o.Position.Scale,
					},
					Data: o.Data,
				}
			}
		}
	}

	if len(e.PreviewViewers) > 0 {
		response.PreviewViewers = make([]res.ViewerPreviewRes, len(e.PreviewViewers))
		for i, v := range e.PreviewViewers {
			response.PreviewViewers[i] = res.ViewerPreviewRes{
				UserID: v.UserID,
				Avatar: v.Avatar,
				Name:   v.Name,
			}
		}
	}

	return response
}
