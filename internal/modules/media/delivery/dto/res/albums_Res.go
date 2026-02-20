package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

type AlbumPrivacyRes struct {
	Level     string   `json:"level"`
	AllowList []string `json:"allow_list"`
	BlockList []string `json:"block_list"`
}

type AlbumReactionStatsRes struct {
	Total int `json:"total"`
	Like  int `json:"like"`
	Love  int `json:"love"`
}

// AlbumRes: Trả về cho client, mọi ObjectID đã được biến thành chuỗi Hex.
type AlbumRes struct {
	ID               string                 `json:"id"`
	UserID           string                 `json:"user_id"`
	GroupID          string                 `json:"group_id,omitempty"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Type             enum.AlbumType         `json:"type"`
	CoverAssetID     string                 `json:"cover_asset_id,omitempty"`
	AssetCount       int                    `json:"asset_count"`
	LastAssetAddedAt *time.Time             `json:"last_asset_added_at,omitempty"`
	Reactions        *AlbumReactionStatsRes `json:"reactions"`
	CommentCount     int                    `json:"comment_count"`
	Privacy          *AlbumPrivacyRes       `json:"privacy"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty"`
}
