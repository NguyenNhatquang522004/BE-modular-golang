package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- REQUEST DTO ---

type AlbumPrivacyReq struct {
	Level     string   `json:"level" validate:"required"`
	AllowList []string `json:"allow_list,omitempty"`
	BlockList []string `json:"block_list,omitempty"`
}

type AlbumReactionStatsReq struct {
	Total int `json:"total"`
	Like  int `json:"like"`
	Love  int `json:"love"`
}

// AlbumReq: Ánh xạ đủ 100% các trường của Entity.
type AlbumReq struct {
	ID               string                 `json:"id,omitempty"` // Dùng string để client gửi hex string
	UserID           string                 `json:"user_id" validate:"required,uuid"`
	GroupID          string                 `json:"group_id,omitempty"`
	Title            string                 `json:"title" validate:"required"`
	Description      string                 `json:"description"`
	Type             enum.AlbumType         `json:"type" validate:"required"`
	CoverAssetID     string                 `json:"cover_asset_id,omitempty"`
	AssetCount       int                    `json:"asset_count"`
	LastAssetAddedAt *time.Time             `json:"last_asset_added_at,omitempty"`
	Reactions        *AlbumReactionStatsReq `json:"reactions"`
	CommentCount     int                    `json:"comment_count"`
	Privacy          *AlbumPrivacyReq       `json:"privacy" validate:"required"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty"`
}

type UpdateAlbumReq struct {
	Title            *string                `json:"title,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Type             *enum.AlbumType        `json:"type,omitempty"`
	CoverAssetID     *string                `json:"cover_asset_id,omitempty"`
	AssetCount       *int                   `json:"asset_count,omitempty"`
	LastAssetAddedAt *time.Time             `json:"last_asset_added_at,omitempty"`
	Reactions        *AlbumReactionStatsReq `json:"reactions,omitempty"`
	CommentCount     *int                   `json:"comment_count,omitempty"`
	Privacy          *AlbumPrivacyReq       `json:"privacy,omitempty"`
}