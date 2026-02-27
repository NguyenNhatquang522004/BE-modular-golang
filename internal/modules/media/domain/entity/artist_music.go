package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionArtist = "Artists"
)

// Artist đại diện cho một nghệ sĩ/nhạc sĩ trên nền tảng mạng xã hội.
type Artist struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. BASIC INFO (Thông tin cơ bản)
	// Index: Text Index { name: "text" } -> Tìm kiếm nghệ sĩ
	Name string `bson:"name" json:"name"`

	// Index: Unique Index { slug: 1 } -> Tối ưu SEO và Routing (VD: "son-tung-m-tp")
	Slug string `bson:"slug" json:"slug"`

	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	AvatarURL string `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	CoverURL  string `bson:"cover_url,omitempty" json:"cover_url,omitempty"` // Ảnh bìa profile

	// 2. VERIFICATION & PLATFORM INTEGRATION
	// Tích xanh chứng nhận nghệ sĩ chính chủ
	IsVerified bool `bson:"is_verified" json:"is_verified"`

	// Liên kết với tài khoản mạng xã hội (Nếu nghệ sĩ này tạo acc trên app của bạn)
	// Index: { user_id: 1 } (Sparse, Unique)
	UserID string `bson:"user_id,omitempty" json:"user_id,omitempty"`

	// 3. SOCIAL LINKS (Mạng xã hội khác)
	SocialLinks SocialLinks `bson:"social_links,omitempty" json:"social_links,omitempty"`

	// 4. METRICS (Chỉ số thống kê - Cần thiết cho Ranking & Profile)
	// Việc lưu metrics ở đây giúp truy vấn nhanh thay vì count() liên tục
	FollowerCount int `bson:"follower_count" json:"follower_count"`
	TotalStreams  int `bson:"total_streams" json:"total_streams"` // Tổng số lần nhạc của họ được dùng

	// 5. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// SocialLinks chứa thông tin các nền tảng khác của nghệ sĩ
type SocialLinks struct {
	Spotify   string `bson:"spotify,omitempty" json:"spotify,omitempty"`
	Youtube   string `bson:"youtube,omitempty" json:"youtube,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Facebook  string `bson:"facebook,omitempty" json:"facebook,omitempty"`
	Website   string `bson:"website,omitempty" json:"website,omitempty"`
}

func (Artist) CollectionName() string {
	return CollectionArtist
}
