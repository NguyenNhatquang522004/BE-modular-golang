package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type CreateAlbumRequest struct {
	// Define fields for creating an album here
	*AlbumReq
	ItemMediaAssetsIDs []string `json:"item_media_assets_ids"`
}

type UpdateAlbumRequest struct {
	AlbumID string `json:"album_id"`
	// Define fields for updating an album here
	*UpdateAlbumReq
}

type DeleteAlbumRequest struct {
	AlbumID string `json:"album_id"`
}

type ReactAlbumRequest struct {
	UserID     string              `json:"user_id"`
	AlbumID    string              `json:"album_id"`
	AssetCount int                 `json:"asset_count"`
	Like       int                 `json:"like"`
	Love       int                 `json:"love"`
	Haha       int                 `json:"haha"`
	Wow        int                 `json:"wow"`
	Sad        int                 `json:"sad"`
	Angry      int                 `json:"angry"`
	EventType  constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
	// Define fields for reacting to an album here
}
type CreateStoryRequest struct {
	*StoryReq
}

type UpdateStoryRequest struct {
	// Define fields for updating a story here
	StoryId string `json:"story_id"`
	UserId  string `json:"user_id"`
	*UpdateStoryReq
}

type DeleteStoryRequest struct {
	StoryId string `json:"story_id"`
	UserId  string `json:"user_id"`
}

type StoryStatsRequest struct {
	UserID          string                           `json:"user_id"`
	StoryID         string                           `json:"story_id"`
	Views           int                              `json:"views"`
	Like            int                              `json:"like"`
	Love            int                              `json:"love"`
	Haha            int                              `json:"haha"`
	Wow             int                              `json:"wow"`
	Sad             int                              `json:"sad"`
	Angry           int                              `json:"angry"`
	ReplyCount      int                              `json:"reply_count"`
	ViewsCount      int                              `json:"views_count"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`
	PollOptionIndex *int                             `json:"poll_option_index,omitempty"` // Chỉ có khi InteractionType là
	Content         string                           `json:"content,omitempty"`           // Chỉ có khi InteractionType là reaction hoặc comment
	EventType       constants.EventType              `json:"event_type"`                  // "increment" hoặc "decrement"
}

type RelyStoryRequest struct {
	// Define fields for relying to a story here
	StoryId         string                           `json:"story_id"`
	UserId          string                           `json:"user_id"`
	Avatar          string                           `json:"avatar"`
	Name            string                           `json:"name"`
	ViewedAt        time.Time                        `json:"viewed_at"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode         `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                             `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType              `json:"event_type"`        // "view" hoặc "unview"
	Content         string                           `json:"content"`
}

type CreateReelRequest struct {
	// Define fields for creating a reel here
	*ReelReq
}

type UpdateReelRequest struct {
	ReelID string `json:"reel_id"`
	UserID string `json:"user_id"`
	// Define fields for updating a reel here
	*UpdateReelReq
}

type DeleteReelRequest struct {
	ReelID string `json:"reel_id"`
	UserID string `json:"user_id"`
}
type ReactReelRequest struct {
	ReelID       string                     `json:"reel_id"`
	UserID       string                     `json:"user_id"`
	Total        int                        `json:"total"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt    time.Time                  `json:"created_at"`
	EventType    constants.EventType        `json:"event_type"` // "view" hoặc "unview"

}

type ReactCounterReelRequest struct {
	ReelID    string              `json:"reel_id"`
	UserID    string              `json:"user_id"`
	Comments  int                 `json:"comments"`
	Saves     int                 `json:"saves"`
	Shares    int                 `json:"shares"`
	EventType constants.EventType `json:"event_type"` // "view" hoặc "unview"

}
type ShareReelRequest struct {
	ReelID string `json:"reel_id"`
	UserID string `json:"user_id"`
	*ReelReq
}

type CreateLiveStreamRequest struct {
	// Define fields for creating a live stream here
	*LiveSessionReq
}

type UpdateLiveStreamRequest struct {
	// Define fields for updating a live stream here
	LiveSessionID string `json:"live_session_id"`
	*UpdateLiveSessionReq
}
type DeleteLiveStreamRequest struct {
	LiveSessionID string `json:"live_session_id"`
	UserID        string `json:"user_id"`
}
type ReactLiveStreamRequest struct {
	LiveSessionID string                     `json:"live_session_id"`
	UserID        string                     `json:"user_id"`
	Total         int                        `json:"total"`
	TargetType    sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode  sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt     time.Time                  `json:"created_at"`
	EventType     constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}

type StartStopVideoLiveStreamRequest struct {
	LiveSessionID string              `json:"live_session_id"`
	SegmentLen    int                 `json:"segment_len"`
	OwnerID       string              `json:"owner_id"`
	Name          string              `json:"name"`
	EventType     constants.EventType `json:"event_type"`
}
type CommentLiveStreamRequest struct {
	*LiveCommentReq
	EventType constants.EventType `json:"event_type"`
}
type CounterLiveStreamRequest struct {
	LiveSessionID string              `json:"live_session_id"`
	Comments      int                 `json:"comments"`
	Views         int                 `json:"views"`
	EventType     constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}
