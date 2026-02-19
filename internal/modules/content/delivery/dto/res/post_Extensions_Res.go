package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

type PostExtensionRes struct {
	ID     string `json:"id"`
	PostID string `json:"post_id"`

	ShareData      *ShareDataRes      `json:"share_data,omitempty"`
	BackgroundData *BackgroundDataRes `json:"background_data,omitempty"`
	QnAData        *QnADataRes        `json:"qna_data,omitempty"`
	ActivityData   *ActivityDataRes   `json:"activity_data,omitempty"`
	LocationDetail *LocationDetailRes `json:"location_detail,omitempty"`
}

// --- Nested Structs ---

type ShareDataRes struct {
	OriginalPostID string           `json:"original_post_id"`
	ParentPostID   string           `json:"parent_post_id"`
	Snapshot       ShareSnapshotRes `json:"snapshot"`
}

type ShareSnapshotRes struct {
	AuthorID       string    `json:"author_id"`
	AuthorName     string    `json:"author_name"`
	AuthorAvatar   string    `json:"author_avatar"`
	ContentExcerpt string    `json:"content_excerpt"`
	MediaThumb     string    `json:"media_thumb,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type BackgroundDataRes struct {
	ThemeID   string `json:"theme_id"`
	TextColor string `json:"text_color"`
}

type QnADataRes struct {
	Question   string `json:"question"`
	ButtonText string `json:"button_text"`
}

type ActivityDataRes struct {
	Type       enum.ActivityType `json:"type"`
	ObjectID   string            `json:"object_id,omitempty"`
	ObjectName string            `json:"object_name"`
}

type LocationDetailRes struct {
	Type        string    `json:"type"` // "Point"
	Coordinates []float64 `json:"coordinates"`
	Address     string    `json:"address"`
	MapURL      string    `json:"map_url,omitempty"`
}