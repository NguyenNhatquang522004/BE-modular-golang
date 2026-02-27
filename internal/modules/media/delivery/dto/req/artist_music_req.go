package req

import "time"

// SocialLinksReq chứa thông tin các nền tảng khác của nghệ sĩ
type SocialLinksReq struct {
	Spotify   string `json:"spotify,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Website   string `json:"website,omitempty"`
}

// ArtistReq đại diện cho payload request từ Client
type ArtistReq struct {
	ID            string         `json:"id,omitempty"` // Nên là string để nhận Hex string từ client
	Name          string         `json:"name" binding:"required"`
	Slug          string         `json:"slug" binding:"required"`
	Bio           string         `json:"bio,omitempty"`
	AvatarURL     string         `json:"avatar_url,omitempty"`
	CoverURL      string         `json:"cover_url,omitempty"`
	IsVerified    bool           `json:"is_verified"`
	UserID        string         `json:"user_id,omitempty"`
	SocialLinks   SocialLinksReq `json:"social_links,omitempty"`
	FollowerCount int            `json:"follower_count"`
	TotalStreams  int            `json:"total_streams"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at,omitempty"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
}