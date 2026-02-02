package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
    CollectionMusicLibrary = "MusicLibrary"
)
type MusicLibrary struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. BASIC INFO
	// Index: Text Index { title: "text", artist: "text", lyrics_snippet: "text" }
	// Mục đích: Tìm kiếm bài hát
	Title  string `bson:"title" json:"title"`
	Artist string `bson:"artist" json:"artist"`
	Album  string `bson:"album,omitempty" json:"album,omitempty"`

	// 2. MEDIA FILES (SeaweedFS / CDN)
	CoverURL string `bson:"cover_url" json:"cover_url"`

	// File nhạc đã cắt sẵn (15s/30s/60s) tối ưu cho streaming
	StreamURL string `bson:"stream_url" json:"stream_url"`

	Duration int `bson:"duration" json:"duration"` // Seconds (VD: 60)

	// 3. SEARCH & DISCOVERY
	// Đoạn lời nổi bật (Hook) để user tìm kiếm bằng lời bài hát
	LyricsSnippet string `bson:"lyrics_snippet,omitempty" json:"lyrics_snippet,omitempty"`

	// Index: Multikey { genres: 1 } -> Lọc nhạc theo thể loại
	Genres []enum.MusicGenre `bson:"genre" json:"genre"`

	// 4. LEGAL & COPYRIGHT
	// Hệ thống cần check field này trước khi cho user dùng nhạc
	CopyrightInfo CopyrightInfo `bson:"copyright_info" json:"copyright_info"`

	// 5. RANKING METRICS
	// Số lần bài hát được sử dụng trong Reels/Stories
	// Index: { usage_count: -1 } -> Trending Music
	UsageCount int `bson:"usage_count" json:"usage_count"`

	// 6. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- COPYRIGHT INFO ---
type CopyrightInfo struct {
	// Nhà cung cấp bản quyền (VD: "Sony Music", "Universal")
	// Dùng String vì danh sách này rất nhiều và biến động
	Provider string `bson:"provider" json:"provider"`

	// Danh sách mã quốc gia ISO (VD: ["VN", "US"])
	// Nếu mảng rỗng hoặc nil -> Coi như Global (toàn cầu)
	AllowedRegions []string `bson:"allowed_regions,omitempty" json:"allowed_regions,omitempty"`
}
func (MusicLibrary) CollectionName() string {
    return CollectionMusicLibrary
}