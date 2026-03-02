package mediaEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

type ReplyStoryPayload struct {
	StoryID string `json:"story_id"`
	ReplyID int    `json:"reply_id"`
}

type CoutnerLiveStreamPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	Comments      int                 `json:"comments"`
	Views         int                 `json:"views"`
	EventType     constants.EventType `json:"event_type"` // "increment" hoặc "decrement"
}
type CommentLiveStreamPayload struct {
	StreamID      gocql.UUID               `json:"stream_id" validate:"required"`
	CreatedAt     time.Time                `json:"created_at"`
	CommentID     gocql.UUID               `json:"comment_id"`
	UserID        gocql.UUID               `json:"user_id" validate:"required"`
	UserNickname  string                   `json:"user_nickname" validate:"required"`
	UserAvatarURL string                   `json:"user_avatar_url"`
	UserBadges    []*sharedEnums.UserBadge `json:"user_badges"` // Sử dụng enum ở DTO
	Content       string                   `json:"content" validate:"required"`
	IsPinned      bool                     `json:"is_pinned"`
	EventType     constants.EventType      `json:"event_type"`
}
type StartStopVideoLiveStreamPayload struct {
	LiveSessionID string              `json:"live_session_id"`
	SegmentLen    int                 `json:"segment_len"`
	OwnerID       string              `json:"owner_id"`
	Name          string              `json:"name"`
	EventType     constants.EventType `json:"event_type"`
}

type ReactLiveStreamPayload struct {
	LiveSessionID string                     `json:"live_session_id"`
	UserID        string                     `json:"user_id"`
	Total         int                        `json:"total"`
	TargetType    sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode  sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt     time.Time                  `json:"created_at"`
	EventType     constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}
type ReactCounterReelPayload struct {
	ReelID    string              `json:"reel_id"`
	UserID    string              `json:"user_id"`
	Comments  int                 `json:"comments"`
	Saves     int                 `json:"saves"`
	Shares    int                 `json:"shares"`
	EventType constants.EventType `json:"event_type"` // "view" hoặc "unview"
}

type ReactReelPayload struct {
	ReelID       string                     `json:"reel_id"`
	UserID       string                     `json:"user_id"`
	Total        int                        `json:"total"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"` // "❤️", "😂" hoặc ID sticker
	CreatedAt    time.Time                  `json:"created_at"`
	EventType    constants.EventType        `json:"event_type"` // "view" hoặc "unview"
}
type ReactStoryPayload struct {
	StoryId         string                           `json:"story_id"`
	UserId          string                           `json:"user_id"`
	Avatar          string                           `json:"avatar"`
	Name            string                           `json:"name"`
	ViewedAt        time.Time                        `json:"viewed_at"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode         `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                             `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType              `json:"event_type"`        // "view" hoặc "unview"
	Content         string                           ` json:"content"`
}

type ViewStoryPayload struct {
	StoryId         string                           `json:"story_id"`
	UserId          string                           `json:"user_id"`
	Avatar          string                           `json:"avatar"`
	Name            string                           `json:"name"`
	ViewedAt        time.Time                        `json:"viewed_at"`
	InteractionType sharedEnums.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode         `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                             `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType              `json:"event_type"`        // "view" hoặc "unview"
	Content         string                           ` json:"content"`
}

type ReactAlbumPayload struct {
	AlbumID        string                   `json:"album_id"`
	TotalReactions int                      `json:"total_reactions"`
	ReactionType   sharedEnums.ReactionCode `json:"reaction_type"`
	EventType      constants.EventType      `json:"event_type"` // "add" hoặc "remove"
}
