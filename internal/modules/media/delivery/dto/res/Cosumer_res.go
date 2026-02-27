package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

type FailedConsumerReactAlbumResponse struct {
	Total        int                      `json:"total"`
	ReactionCode sharedEnums.ReactionCode `json:"reaction_code"`
}

type FailedConsumerReactStoryResponse struct {
	StoryId         string                    `json:"story_id"`
	UserId          string                    `json:"user_id"`
	Avatar          string                    `json:"avatar"`
	Name            string                    `json:"name"`
	ViewedAt        time.Time                 `json:"viewed_at"`
	InteractionType enum.StoryInteractionType `json:"interaction_type"`  // "view", "reaction", hoặc "poll_vote"
	ReactionCode    sharedEnums.ReactionCode  `json:"reaction_code"`     // "❤️", "😂" hoặc ID sticker
	PollOptionIndex *int                      `json:"poll_option_index"` // nil nếu không vote, 0 hoặc 1 nếu có vote
	EventType       constants.EventType       `json:"event_type"`        // "view" hoặc "unview"
	ErrorMessage    string                    `json:"error_message"`
}
