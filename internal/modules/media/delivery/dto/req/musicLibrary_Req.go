package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB-STRUCTS ---

type CopyrightInfoReq struct {
	Provider       string   `json:"provider" validate:"required"`
	AllowedRegions []string `json:"allowed_regions,omitempty"`
}

// --- MAIN REQUEST DTO ---

// MusicLibraryReq: Ánh xạ 100% các trường của Entity.
type MusicLibraryReq struct {
	ID            string                   `json:"id,omitempty"` // Trình bày dưới dạng string để nhận Hex ObjectID
	Title         string                   `json:"title" validate:"required"`
	ArtistID      string                   `json:"artist_id" validate:"required"`
	Album         string                   `json:"album,omitempty"`
	CoverURL      string                   `json:"cover_url"`
	StreamURL     string                   `json:"stream_url" validate:"required"`
	Duration      int                      `json:"duration"`
	LyricsSnippet string                   `json:"lyrics_snippet,omitempty"`
	Genres        []sharedEnums.MusicGenre `json:"genre"` // json tag giữ là "genre" giống entity
	CopyrightInfo CopyrightInfoReq         `json:"copyright_info"`
	UsageCount    int                      `json:"usage_count"`
	CreatedAt     *time.Time               `json:"created_at,omitempty"`
	UpdatedAt     *time.Time               `json:"updated_at,omitempty"`
	DeletedAt     *time.Time               `json:"deleted_at,omitempty"`
}

// UpdateMusicLibraryReq: Áp dụng 100% pointer để hỗ trợ Partial Update
type UpdateMusicLibraryReq struct {
	Title         *string                  `json:"title,omitempty"`
	ArtistID      *string                  `json:"artist_id,omitempty"`
	Album         *string                  `json:"album,omitempty"`
	CoverURL      *string                  `json:"cover_url,omitempty"`
	StreamURL     *string                  `json:"stream_url,omitempty"`
	Duration      *int                     `json:"duration,omitempty"`
	LyricsSnippet *string                  `json:"lyrics_snippet,omitempty"`
	Genres        []sharedEnums.MusicGenre `json:"genre,omitempty"`
	CopyrightInfo *CopyrightInfoReq        `json:"copyright_info,omitempty"`
	UsageCount    *int                     `json:"usage_count,omitempty"`
	DeletedAt     *time.Time               `json:"deleted_at,omitempty"`
}
