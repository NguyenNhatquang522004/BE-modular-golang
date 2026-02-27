package utils

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

func CreateReactAlbumRequestToAlbumReactionStats(data *entity.Album, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactAlbumRequest to AlbumReactionStats here
	data.Reactions.Total = +1
	switch enum {
	case sharedEnums.ReactionLike:
		data.Reactions.Like = +1
	case sharedEnums.ReactionLove:
		data.Reactions.Love = +1
	case sharedEnums.ReactionHaha:
		data.Reactions.Haha = +1
	case sharedEnums.ReactionWow:
		data.Reactions.Wow = +1
	case sharedEnums.ReactionSad:
		data.Reactions.Sad = +1
	case sharedEnums.ReactionAngry:
		data.Reactions.Angry = +1
	default:
		// Handle unknown reaction type
	}
}

func DeleteReactAlbumRequestToAlbumReactionStats(data *entity.Album, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactAlbumRequest to AlbumReactionStats here
	data.Reactions.Total = -1
	switch enum {
	case sharedEnums.ReactionLike:
		data.Reactions.Like = -1
	case sharedEnums.ReactionLove:
		data.Reactions.Love = -1
	case sharedEnums.ReactionHaha:
		data.Reactions.Haha = -1
	case sharedEnums.ReactionWow:
		data.Reactions.Wow = -1
	case sharedEnums.ReactionSad:
		data.Reactions.Sad = -1
	case sharedEnums.ReactionAngry:
		data.Reactions.Angry = -1
	default:
		// Handle unknown reaction type
	}
}
func CreateReactStoryRequestToStoryReactionStats(data *entity.Story, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactStoryRequest to StoryReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Likes = +1
	case sharedEnums.ReactionLove:
		data.Stats.Love = +1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = +1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = +1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = +1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = +1
	default:
		// Handle unknown reaction type
	}
}

func DeleteReactStoryRequestToStoryReactionStats(data *entity.Story, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactStoryRequest to StoryReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Likes = -1
	case sharedEnums.ReactionLove:
		data.Stats.Love = -1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = -1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = -1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = -1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = -1
	default:
		// Handle unknown reaction type
	}
}
func CreateReactReelRequestToReelReactionStats(data *entity.Reel, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactReelRequest to ReelReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Like = +1
	case sharedEnums.ReactionLove:
		data.Stats.Love = +1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = +1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = +1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = +1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = +1
	default:
		// Handle unknown reaction type
	}
}

func DeleteReactReelRequestToReelReactionStats(data *entity.Reel, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactReelRequest to ReelReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Like = -1
	case sharedEnums.ReactionLove:
		data.Stats.Love = -1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = -1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = -1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = -1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = -1
	default:
		// Handle unknown reaction type
	}
}
func CreateReactLiveSessionRequestToLiveSessionReactionStats(data *entity.LiveSession, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactLiveSessionRequest to LiveSessionReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Like = +1
	case sharedEnums.ReactionLove:
		data.Stats.Love = +1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = +1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = +1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = +1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = +1
	default:
		// Handle unknown reaction type
	}
}

func DeleteReactLiveSessionRequestToLiveSessionReactionStats(data *entity.LiveSession, enum sharedEnums.ReactionCode) {
	// Implement the logic to map ReactLiveSessionRequest to LiveSessionReactionStats here
	switch enum {
	case sharedEnums.ReactionLike:
		data.Stats.Like = -1
	case sharedEnums.ReactionLove:
		data.Stats.Love = -1
	case sharedEnums.ReactionHaha:
		data.Stats.Haha = -1
	case sharedEnums.ReactionWow:
		data.Stats.Wow = -1
	case sharedEnums.ReactionSad:
		data.Stats.Sad = -1
	case sharedEnums.ReactionAngry:
		data.Stats.Angry = -1
	default:
		// Handle unknown reaction type
	}
}
