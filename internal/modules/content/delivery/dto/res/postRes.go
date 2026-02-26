	package res

	import (
		"time"

		"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
	)

	// --- NESTED DTOs ---
	type PostContextRes struct {
		Type     enum.ContextType `json:"type"`
		TargetID string           `json:"target_id,omitempty"`
	}

	type PostSummaryRes struct {
		FeelingIcon       string `json:"feeling_icon,omitempty"`
		FeelingName       string `json:"feeling_name,omitempty"`
		LocationName      string `json:"location_name,omitempty"`
		HasMedia          bool   `json:"has_media"`
		MediaCount        int    `json:"media_count"`
		ThumbnailURL      string `json:"thumbnail_url,omitempty"`
		BackgroundThemeID string `json:"background_theme_id,omitempty"`
	}

	type PostPrivacyRes struct {
		Scope        enum.PrivacyScope `json:"scope"`
		AllowComment bool              `json:"allow_comment"`
		AllowShare   bool              `json:"allow_share"`
	}

	type PostStatsRes struct {
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
	type PostRes struct {
		ID          string          `json:"id"` // Trả về dạng string cho client
		UserID      string          `json:"user_id"`
		Type        enum.PostType   `json:"type"`
		Context     *PostContextRes `json:"context,omitempty"`
		Content     string          `json:"content"`
		Slug        string          `json:"slug"`
		Summary     *PostSummaryRes `json:"summary,omitempty"`
		Privacy     PostPrivacyRes  `json:"privacy"`
		Status      enum.PostStatus `json:"status"`
		IsPinned    bool            `json:"is_pinned"`
		IsEdited    bool            `json:"is_edited"`
		Stats       PostStatsRes    `json:"stats"`
		Hashtags    []string        `json:"hashtags,omitempty"`
		Mentions    []string        `json:"mentions,omitempty"`
		PublishedAt *time.Time      `json:"published_at,omitempty"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
		DeletedAt   *time.Time      `json:"deleted_at,omitempty"`
	}
