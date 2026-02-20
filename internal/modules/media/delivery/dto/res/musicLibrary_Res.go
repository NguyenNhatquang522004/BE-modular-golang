package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type CopyrightInfoRes struct {
	Provider       string   `json:"provider"`
	AllowedRegions []string `json:"allowed_regions,omitempty"`
}

// --- MAIN RESPONSE DTO ---

// MusicLibraryRes: Cấu trúc trả về cho Client
type MusicLibraryRes struct {
	ID            string            `json:"id"` // Trả về Hex string an toàn
	Title         string            `json:"title"`
	Artist        string            `json:"artist"`
	Album         string            `json:"album,omitempty"`
	CoverURL      string            `json:"cover_url"`
	StreamURL     string            `json:"stream_url"`
	Duration      int               `json:"duration"`
	LyricsSnippet string            `json:"lyrics_snippet,omitempty"`
	Genres        []enum.MusicGenre `json:"genre"`
	CopyrightInfo CopyrightInfoRes  `json:"copyright_info"`
	UsageCount    int               `json:"usage_count"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     *time.Time        `json:"deleted_at,omitempty"`
}