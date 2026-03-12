package mapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToMusicLibraryEntity(req *mediaEvent.CreateMusicLibraryPayload) (*entity.MusicLibrary, error) {
	if req == nil {
		return nil, errors.New("create dto cannot be nil")
	}

	// 1. Chuyển đổi ArtistID từ chuỗi hex (string) sang primitive.ObjectID
	artistID, err := primitive.ObjectIDFromHex(req.ArtistID)
	if err != nil {
		return nil, fmt.Errorf("invalid artist_id format: %w", err)
	}

	now := time.Now()

	// 2. Mapping dữ liệu và thiết lập các trường hệ thống
	return &entity.MusicLibrary{
		// ID tự động sinh ra cho bản ghi mới
		ID: primitive.NewObjectID(),

		ArtistID:      artistID,
		Title:         req.Title,
		Album:         req.Album,
		CoverURL:      req.CoverURL,
		StreamURL:     req.StreamURL,
		Duration:      req.Duration,
		LyricsSnippet: req.LyricsSnippet,
		Genres:        req.Genres,

		CopyrightInfo: entity.CopyrightInfo{
			Provider:       req.CopyrightInfo.Provider,
			AllowedRegions: req.CopyrightInfo.AllowedRegions,
		},

		// Các trường do hệ thống kiểm soát (System fields)
		UsageCount: 0, // Mặc định bài hát mới chưa có ai sử dụng
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func ApplyUpdateMusicLibrary(ent *entity.MusicLibrary, req *mediaEvent.UpdateMusicLibraryPayload) *entity.MusicLibrary {
	if ent == nil || req == nil {
		return ent
	}
	hasUpdate := false // Cờ đánh dấu xem có thực sự xảy ra cập nhật không

	// Kiểm tra từng con trỏ, nếu khác nil tức là client CÓ GỬI trường này lên để cập nhật
	if req.Title != nil {
		ent.Title = *req.Title
		hasUpdate = true
	}
	if req.Album != nil {
		ent.Album = *req.Album
		hasUpdate = true
	}
	if req.CoverURL != nil {
		ent.CoverURL = *req.CoverURL
		hasUpdate = true
	}
	if req.StreamURL != nil {
		ent.StreamURL = *req.StreamURL
		hasUpdate = true
	}
	if req.Duration != nil {
		ent.Duration = *req.Duration
		hasUpdate = true
	}
	if req.LyricsSnippet != nil {
		ent.LyricsSnippet = *req.LyricsSnippet
		hasUpdate = true
	}
	if req.Genres != nil {
		ent.Genres = *req.Genres
		hasUpdate = true
	}

	// Xử lý struct lồng nhau (Nested Struct)
	if req.CopyrightInfo != nil {
		if req.CopyrightInfo.Provider != nil {
			ent.CopyrightInfo.Provider = *req.CopyrightInfo.Provider
			hasUpdate = true
		}
		if req.CopyrightInfo.AllowedRegions != nil {
			ent.CopyrightInfo.AllowedRegions = *req.CopyrightInfo.AllowedRegions
			hasUpdate = true
		}
	}

	// Chỉ cập nhật UpdatedAt nếu thực sự có field nào đó bị thay đổi
	if hasUpdate {
		ent.UpdatedAt = time.Now()
	}

	return ent
}
