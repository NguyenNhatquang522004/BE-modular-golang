package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
    CollectionAlbums = "Albums"
)
type Album struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. OWNERSHIP
	// UserID từ Postgres (UUID) -> Lưu String
	// Index: { user_id: 1, type: 1 } -> Lấy list album của user
	UserID string `bson:"user_id" json:"user_id"`

	// GroupID từ Mongo. Dùng Pointer để có thể null.
	// Index: { group_id: 1 } -> Lấy list album của nhóm
	GroupID *primitive.ObjectID `bson:"group_id,omitempty" json:"group_id,omitempty"`

	// 2. BASIC INFO
	Title       string         `bson:"title" json:"title"`
	Description string         `bson:"description" json:"description"`
	Type        enum.AlbumType `bson:"type" json:"type"` // 'normal', 'profile'...

	// 3. COVER IMAGE
	// ID của MediaAsset (Mongo ID)
	CoverAssetID primitive.ObjectID `bson:"cover_asset_id,omitempty" json:"cover_asset_id,omitempty"`

	// 4. DENORMALIZATION (UI Optimizations)
	// Tổng số ảnh/video trong album
	AssetCount int `bson:"asset_count" json:"asset_count"`

	// Thời điểm có ảnh mới nhất được thêm vào (Dùng để Sort list album)
	// Index: { user_id: 1, last_asset_added_at: -1 }
	LastAssetAddedAt *time.Time `bson:"last_asset_added_at,omitempty" json:"last_asset_added_at,omitempty"`

	// 5. INTERACTION & PRIVACY
	Reactions    AlbumReactionStats `bson:"reactions" json:"reactions"`
	CommentCount int                `bson:"comment_count" json:"comment_count"`
	Privacy      AlbumPrivacy       `bson:"privacy" json:"privacy"`

	// 6. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- PRIVACY ---
type AlbumPrivacy struct {
	// Level: 'public', 'friends', 'only_me', 'custom'
	// Có thể dùng string hoặc tái sử dụng Enum PrivacyScope
	Level string `bson:"level" json:"level"`

	// Danh sách User ID được phép xem (UUID String)
	AllowList []string `bson:"allow_list,omitempty" json:"allow_list,omitempty"`

	// Danh sách User ID bị chặn xem (UUID String)
	BlockList []string `bson:"block_list,omitempty" json:"block_list,omitempty"`
}

// --- REACTIONS ---
// Denormalization: Lưu tổng số reaction để hiển thị ngay bên ngoài album
type AlbumReactionStats struct {
	Total int `bson:"total" json:"total"`
	Like  int `bson:"like" json:"like"`
	Love  int `bson:"love" json:"love"`
}
func (Album) CollectionName() string {
    return CollectionAlbums
}