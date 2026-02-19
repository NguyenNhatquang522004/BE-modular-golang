package res

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"

type PostMediaRes struct {
	ID     string          `json:"id"`
	PostID string          `json:"post_id"`
	Items  []*MediaItemRes `json:"items"`
}

// --- Nested Structs ---

type MediaItemRes struct {
	ID           string           `json:"id"`
	MediaType    enum.MediaType   `json:"media_type"`
	URL          string           `json:"url"`
	ThumbnailURL string           `json:"thumbnail_url"`
	Metadata     MediaMetadataRes `json:"metadata"`
	Order        int              `json:"order"`
	TaggedUsers  []TaggedUserRes  `json:"tagged_users,omitempty"`
}

type MediaMetadataRes struct {
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes"`
	MimeType  string  `json:"mime_type"`
}

type TaggedUserRes struct {
	UserID string  `json:"user_id"`
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}
