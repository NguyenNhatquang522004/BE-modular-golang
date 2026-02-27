package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntityLiveSession: Chuyển từ Req -> Entity để Insert
func ToEntityLiveSession(r *req.LiveSessionReq) (*entity.LiveSession, error) {
	if r == nil {
		return nil, nil
	}

	e := &entity.LiveSession{
		HostUserID:      r.HostUserID,
		Title:           r.Title,
		Description:     r.Description,
		CategoryID:      r.CategoryID,
		Status:          r.Status,
		StreamKey:       r.StreamKey,
		PlaybackURL:     r.PlaybackURL,
		BannedUsers:     r.BannedUsers,
		PinnedCommentID: r.PinnedCommentID,
		StartedAt:       r.StartedAt,
		EndedAt:         r.EndedAt,
		RecordingSetting: entity.RecordingSetting{
			IsRecorded: r.RecordingSetting.IsRecorded,
			ArchiveURL: r.RecordingSetting.ArchiveURL,
		},
		Stats: entity.LiveStats{
			PeakViewers:   r.Stats.PeakViewers,
			TotalViews:    r.Stats.TotalViews,
			TotalComments: r.Stats.TotalComments,
			Like:          r.Stats.Like,
			Love:          r.Stats.Love,
			Haha:          r.Stats.Haha,
			Wow:           r.Stats.Wow,
			Sad:           r.Stats.Sad,
			Angry:         r.Stats.Angry,
		},
	}

	// Xử lý ID
	if r.ID != "" {
		id, err := primitive.ObjectIDFromHex(r.ID)
		if err == nil {
			e.ID = id
		}
	} else {
		e.ID = primitive.NewObjectID()
	}

	// Xử lý Timestamps
	now := time.Now()
	if r.CreatedAt != nil {
		e.CreatedAt = *r.CreatedAt
	} else {
		e.CreatedAt = now
	}

	if r.UpdatedAt != nil {
		e.UpdatedAt = *r.UpdatedAt
	} else {
		e.UpdatedAt = now
	}

	return e, nil
}

// 2. UpdateToEntityLiveSession: Đổ dữ liệu từ UpdateReq vào Entity có sẵn (Partial Update)
func UpdateToEntityLiveSession(r *req.UpdateLiveSessionReq, e *entity.LiveSession) {
	if r == nil || e == nil {
		return
	}

	if r.HostUserID != nil {
		e.HostUserID = *r.HostUserID
	}
	if r.Title != nil {
		e.Title = *r.Title
	}
	if r.Description != nil {
		e.Description = *r.Description
	}
	if r.CategoryID != nil {
		e.CategoryID = *r.CategoryID
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if r.StreamKey != nil {
		e.StreamKey = *r.StreamKey
	}
	if r.PlaybackURL != nil {
		e.PlaybackURL = *r.PlaybackURL
	}
	if r.RecordingSetting != nil {
		e.RecordingSetting.IsRecorded = r.RecordingSetting.IsRecorded
		e.RecordingSetting.ArchiveURL = r.RecordingSetting.ArchiveURL
	}
	if r.BannedUsers != nil {
		e.BannedUsers = r.BannedUsers // Replace slice
	}
	if r.PinnedCommentID != nil {
		e.PinnedCommentID = *r.PinnedCommentID
	}
	if r.StartedAt != nil {
		e.StartedAt = r.StartedAt
	}
	if r.EndedAt != nil {
		e.EndedAt = r.EndedAt
	}
	if r.Stats != nil {
		e.Stats.PeakViewers = r.Stats.PeakViewers
		e.Stats.TotalViews = r.Stats.TotalViews
		e.Stats.Like = r.Stats.Like
		e.Stats.Love = r.Stats.Love
		e.Stats.Haha = r.Stats.Haha
		e.Stats.Wow = r.Stats.Wow
		e.Stats.Sad = r.Stats.Sad
		e.Stats.Angry = r.Stats.Angry
		e.Stats.TotalComments = r.Stats.TotalComments
	}

	e.UpdatedAt = time.Now()
}

// 3. ReqToResLiveSession: Chuyển trực tiếp từ Req -> Res
func ReqToResLiveSession(r *req.LiveSessionReq) *res.LiveSessionRes {
	if r == nil {
		return nil
	}

	response := &res.LiveSessionRes{
		ID:              r.ID,
		HostUserID:      r.HostUserID,
		Title:           r.Title,
		Description:     r.Description,
		CategoryID:      r.CategoryID,
		Status:          r.Status,
		StreamKey:       r.StreamKey,
		PlaybackURL:     r.PlaybackURL,
		BannedUsers:     r.BannedUsers,
		PinnedCommentID: r.PinnedCommentID,
		StartedAt:       r.StartedAt,
		EndedAt:         r.EndedAt,
		RecordingSetting: res.RecordingSettingRes{
			IsRecorded: r.RecordingSetting.IsRecorded,
			ArchiveURL: r.RecordingSetting.ArchiveURL,
		},
		Stats: res.LiveStatsRes{
			PeakViewers:   r.Stats.PeakViewers,
			TotalViews:    r.Stats.TotalViews,
			TotalComments: r.Stats.TotalComments,
			Like:          r.Stats.Like,
			Love:          r.Stats.Love,
			Haha:          r.Stats.Haha,
			Wow:           r.Stats.Wow,
			Sad:           r.Stats.Sad,
			Angry:         r.Stats.Angry,
		},
	}

	now := time.Now()
	if r.CreatedAt != nil {
		response.CreatedAt = *r.CreatedAt
	} else {
		response.CreatedAt = now
	}

	if r.UpdatedAt != nil {
		response.UpdatedAt = *r.UpdatedAt
	} else {
		response.UpdatedAt = now
	}

	return response
}

// 4. EntityToResLiveSession (Bổ sung Best Practice): Chuyển Entity thành Response
func EntityToResLiveSession(e *entity.LiveSession) *res.LiveSessionRes {
	if e == nil {
		return nil
	}

	return &res.LiveSessionRes{
		ID:              e.ID.Hex(),
		HostUserID:      e.HostUserID,
		Title:           e.Title,
		Description:     e.Description,
		CategoryID:      e.CategoryID,
		Status:          e.Status,
		StreamKey:       e.StreamKey,
		PlaybackURL:     e.PlaybackURL,
		BannedUsers:     e.BannedUsers,
		PinnedCommentID: e.PinnedCommentID,
		StartedAt:       e.StartedAt,
		EndedAt:         e.EndedAt,
		RecordingSetting: res.RecordingSettingRes{
			IsRecorded: e.RecordingSetting.IsRecorded,
			ArchiveURL: e.RecordingSetting.ArchiveURL,
		},
		Stats: res.LiveStatsRes{
			PeakViewers:   e.Stats.PeakViewers,
			TotalViews:    e.Stats.TotalViews,
			TotalComments: e.Stats.TotalComments,
			Like:          e.Stats.Like,
			Love:          e.Stats.Love,
			Haha:          e.Stats.Haha,
			Wow:           e.Stats.Wow,
			Sad:           e.Stats.Sad,
			Angry:         e.Stats.Angry,
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
