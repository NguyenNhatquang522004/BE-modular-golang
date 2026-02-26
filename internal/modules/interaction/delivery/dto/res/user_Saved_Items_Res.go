package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type UserSavedItemRes struct {
	ID             string                       `json:"id"`
	UserID         string                       `json:"user_id"`
	TargetID       string                       `json:"target_id"`
	TargetType     *sharedEnums.SavedTargetType `json:"target_type"`
	Snapshot       *SavedItemSnapshotRes        `json:"snapshot"`
	CollectionName string                       `json:"collection_name"`
	CreatedAt      time.Time                    `json:"created_at"`
}

type SavedItemSnapshotRes struct {
	AuthorName     string `json:"author_name"`
	ContentPreview string `json:"content_preview"`
	ThumbnailURL   string `json:"thumbnail_url"`
}
