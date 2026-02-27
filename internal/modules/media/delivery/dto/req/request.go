package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
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
	AlbumID        string                   `json:"album_id"`
	TotalReactions int                      `json:"total_reactions"`
	ReactionType   sharedEnums.ReactionCode `json:"reaction_type"`
	EventType      constants.EventType      `json:"event_type"` // "add" hoặc "remove"

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

type ViewCountStoryRequest struct {
	StoryId         string                    `json:"story_id"`
	UserId          string                    `json:"user_id"`
	Avatar          string                    `json:"avatar"`
	Name            string                    `json:"name"`
	ViewedAt        time.Time                 `json:"viewed_at"`
	InteractionType enum.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode  `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                      `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType       `json:"event_type"`        // "view" hoặc "unview"
}
type ReactStoryRequest struct {
	StoryId         string                    `json:"story_id"`
	UserId          string                    `json:"user_id"`
	Avatar          string                    `json:"avatar"`
	Name            string                    `json:"name"`
	ViewedAt        time.Time                 `json:"viewed_at"`
	InteractionType enum.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode  `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                      `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType       `json:"event_type"`        // "view" hoặc "unview"
}
