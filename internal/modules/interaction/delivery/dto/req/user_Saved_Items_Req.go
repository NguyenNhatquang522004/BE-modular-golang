package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type UserSavedItemReq struct {
	ID             string                       `json:"id"`                              // Rỗng khi Create, dùng khi Update
	UserID         string                       `json:"user_id" binding:"required,uuid"` // Đảm bảo đúng định dạng UUID từ Postgres
	TargetID       string                       `json:"target_id" binding:"required"`    // Chuỗi Hex của ObjectID
	TargetType     *sharedEnums.SavedTargetType `json:"target_type" binding:"required"`
	Snapshot       *SavedItemSnapshotReq        `json:"snapshot" binding:"required"`
	CollectionName string                       `json:"collection_name"`
	CreatedAt      time.Time                    `json:"created_at"`
}

type SavedItemSnapshotReq struct {
	AuthorName     string `json:"author_name" binding:"required"`
	ContentPreview string `json:"content_preview"`
	ThumbnailURL   string `json:"thumbnail_url" binding:"omitempty,url"` // Kiểm tra URL hợp lệ nếu có truyền
}
