package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. ToEntityMusicLibrary: Chuyển từ Req -> Entity để Insert
func ToEntityMusicLibrary(r *req.MusicLibraryReq) (*entity.MusicLibrary, error) {
	if r == nil {
		return nil, nil
	}

	e := &entity.MusicLibrary{
		Title:         r.Title,
		Artist:        r.Artist,
		Album:         r.Album,
		CoverURL:      r.CoverURL,
		StreamURL:     r.StreamURL,
		Duration:      r.Duration,
		LyricsSnippet: r.LyricsSnippet,
		Genres:        r.Genres,
		UsageCount:    r.UsageCount,
		CopyrightInfo: entity.CopyrightInfo{
			Provider:       r.CopyrightInfo.Provider,
			AllowedRegions: r.CopyrightInfo.AllowedRegions,
		},
		DeletedAt: r.DeletedAt,
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

// 2. UpdateToEntityMusicLibrary: Cập nhật Partial Update (chỉ update các trường được gửi lên)
func UpdateToEntityMusicLibrary(r *req.UpdateMusicLibraryReq, e *entity.MusicLibrary) {
	if r == nil || e == nil {
		return
	}

	if r.Title != nil {
		e.Title = *r.Title
	}
	if r.Artist != nil {
		e.Artist = *r.Artist
	}
	if r.Album != nil {
		e.Album = *r.Album
	}
	if r.CoverURL != nil {
		e.CoverURL = *r.CoverURL
	}
	if r.StreamURL != nil {
		e.StreamURL = *r.StreamURL
	}
	if r.Duration != nil {
		e.Duration = *r.Duration
	}
	if r.LyricsSnippet != nil {
		e.LyricsSnippet = *r.LyricsSnippet
	}
	if r.Genres != nil {
		e.Genres = r.Genres
	}
	if r.CopyrightInfo != nil {
		e.CopyrightInfo.Provider = r.CopyrightInfo.Provider
		e.CopyrightInfo.AllowedRegions = r.CopyrightInfo.AllowedRegions
	}
	if r.UsageCount != nil {
		e.UsageCount = *r.UsageCount
	}
	if r.DeletedAt != nil {
		e.DeletedAt = r.DeletedAt
	}

	e.UpdatedAt = time.Now()
}

// 3. ReqToResMusicLibrary: Chuyển thẳng từ Req -> Res
func ReqToResMusicLibrary(r *req.MusicLibraryReq) *res.MusicLibraryRes {
	if r == nil {
		return nil
	}

	response := &res.MusicLibraryRes{
		ID:            r.ID,
		Title:         r.Title,
		Artist:        r.Artist,
		Album:         r.Album,
		CoverURL:      r.CoverURL,
		StreamURL:     r.StreamURL,
		Duration:      r.Duration,
		LyricsSnippet: r.LyricsSnippet,
		Genres:        r.Genres,
		UsageCount:    r.UsageCount,
		CopyrightInfo: res.CopyrightInfoRes{
			Provider:       r.CopyrightInfo.Provider,
			AllowedRegions: r.CopyrightInfo.AllowedRegions,
		},
		DeletedAt: r.DeletedAt,
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

// 4. EntityToResMusicLibrary: Chuyển đổi dữ liệu truy vấn từ DB -> Chuẩn Response cho Client
func EntityToResMusicLibrary(e *entity.MusicLibrary) *res.MusicLibraryRes {
	if e == nil {
		return nil
	}

	return &res.MusicLibraryRes{
		ID:            e.ID.Hex(),
		Title:         e.Title,
		Artist:        e.Artist,
		Album:         e.Album,
		CoverURL:      e.CoverURL,
		StreamURL:     e.StreamURL,
		Duration:      e.Duration,
		LyricsSnippet: e.LyricsSnippet,
		Genres:        e.Genres,
		UsageCount:    e.UsageCount,
		CopyrightInfo: res.CopyrightInfoRes{
			Provider:       e.CopyrightInfo.Provider,
			AllowedRegions: e.CopyrightInfo.AllowedRegions,
		},
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}
