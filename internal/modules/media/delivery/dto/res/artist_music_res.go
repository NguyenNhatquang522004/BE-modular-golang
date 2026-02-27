package res

import "time"

// SocialLinksRes chứa thông tin các nền tảng khác của nghệ sĩ
type SocialLinksRes struct {
	Spotify   string `json:"spotify,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Website   string `json:"website,omitempty"`
}

// ArtistRes đại diện cho payload response trả về cho Client
type ArtistRes struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	Bio           string         `json:"bio,omitempty"`
	AvatarURL     string         `json:"avatar_url,omitempty"`
	CoverURL      string         `json:"cover_url,omitempty"`
	IsVerified    bool           `json:"is_verified"`
	UserID        string         `json:"user_id,omitempty"`
	SocialLinks   SocialLinksRes `json:"social_links,omitempty"`
	FollowerCount int            `json:"follower_count"`
	TotalStreams  int            `json:"total_streams"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
}