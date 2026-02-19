package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

type PostResponse struct {
	ID     string `json:"id"`      // Đã convert từ ObjectID
	UserID string `json:"user_id"` // UUID

	Type    enum.PostType        `json:"type"`
	Context *PostContextResponse `json:"context,omitempty"`

	Content string               `json:"content"`
	Slug    string               `json:"slug"`
	Summary *PostSummaryResponse `json:"summary,omitempty"`

	Privacy  PostPrivacyResponse `json:"privacy"`
	Status   enum.PostStatus     `json:"status"`
	IsPinned bool                `json:"is_pinned"`
	IsEdited bool                `json:"is_edited"`

	Stats PostStatsResponse `json:"stats"`

	Hashtags []string `json:"hashtags,omitempty"`
	Mentions []string `json:"mentions,omitempty"`

	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// --- SUB-STRUCTS FOR RESPONSE ---

type PostContextResponse struct {
	Type     enum.ContextType `json:"type"`
	TargetID string           `json:"target_id"`
}

type PostSummaryResponse struct {
	FeelingIcon       string `json:"feeling_icon,omitempty"`
	FeelingName       string `json:"feeling_name,omitempty"`
	LocationName      string `json:"location_name,omitempty"`
	HasMedia          bool   `json:"has_media"`
	MediaCount        int    `json:"media_count"`
	ThumbnailURL      string `json:"thumbnail_url,omitempty"`
	BackgroundThemeID string `json:"background_theme_id,omitempty"`
}

type PostPrivacyResponse struct {
	Scope        enum.PrivacyScope `json:"scope"`
	AllowComment bool              `json:"allow_comment"`
	AllowShare   bool              `json:"allow_share"`
}

type PostStatsResponse struct {
	TotalReactions   int      `json:"total_reactions"`
	Comments         int      `json:"comments"`
	Shares           int      `json:"shares"`
	Views            int      `json:"views"`
	TopReactionTypes []string `json:"top_reaction_types"`
}
