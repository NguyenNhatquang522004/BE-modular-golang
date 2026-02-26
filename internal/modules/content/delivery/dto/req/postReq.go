package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// --- NESTED DTOs ---
type PostContextReq struct {
	Type     enum.ContextType `json:"type" validate:"required"`
	TargetID string           `json:"target_id,omitempty"`
}

type PostSummaryReq struct {
	FeelingIcon       string `json:"feeling_icon,omitempty"`
	FeelingName       string `json:"feeling_name,omitempty"`
	LocationName      string `json:"location_name,omitempty"`
	HasMedia          bool   `json:"has_media"`
	MediaCount        int    `json:"media_count"`
	ThumbnailURL      string `json:"thumbnail_url,omitempty"`
	BackgroundThemeID string `json:"background_theme_id,omitempty"`
}

type PostPrivacyReq struct {
	Scope        sharedEnums.PrivacyScope `json:"scope" validate:"required"`
	AllowComment bool                     `json:"allow_comment"`
	AllowShare   bool                     `json:"allow_share"`
}

type PostStatsReq struct {
	TotalReactions   int      `json:"total_reactions"`
	Comments         int      `json:"comments"`
	Shares           int      `json:"shares"`
	Views            int      `json:"views"`
	TopReactionTypes []string `json:"top_reaction_types"`
	Like             int      `json:"like"`
	Love             int      `json:"love"`
	Haha             int      `json:"haha"`
	Wow              int      `json:"wow"`
	Sad              int      `json:"sad"`
	Angry            int      `json:"angry"`
}

// --- MAIN DTO ---
type PostReq struct {
	ID          string          `json:"id,omitempty"` // String thay vì ObjectID cho HTTP Request
	UserID      string          `json:"user_id" validate:"required"`
	Type        enum.PostType   `json:"type" validate:"required"`
	Context     *PostContextReq `json:"context,omitempty"`
	Content     string          `json:"content" validate:"required"`
	Slug        string          `json:"slug"`
	Summary     *PostSummaryReq `json:"summary,omitempty"`
	Privacy     PostPrivacyReq  `json:"privacy"`
	Status      sharedEnums.ProcessingStatus `json:"status" validate:"required"`
	IsPinned    bool            `json:"is_pinned"`
	IsEdited    bool            `json:"is_edited"`
	Stats       PostStatsReq    `json:"stats"`
	Hashtags    []string        `json:"hashtags,omitempty"`
	Mentions    []string        `json:"mentions,omitempty"`
	PublishedAt *time.Time      `json:"published_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty"`
}
