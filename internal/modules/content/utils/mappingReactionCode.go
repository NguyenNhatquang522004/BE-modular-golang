package utils

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

func MappingCreateReactioncode(datacomment *entity.Post, react *sharedEnums.ReactionCode) {
	datacomment.Stats.TotalReactions += 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Stats.Like += 1
	case sharedEnums.ReactionLove:
		datacomment.Stats.Love += 1
	case sharedEnums.ReactionHaha:
		datacomment.Stats.Haha += 1
	case sharedEnums.ReactionSad:
		datacomment.Stats.Sad += 1
	case sharedEnums.ReactionAngry:
		datacomment.Stats.Angry += 1
	case sharedEnums.ReactionWow:
		datacomment.Stats.Wow += 1
	}
}

func MappingDeleteReactioncode(datacomment *entity.Post, react *sharedEnums.ReactionCode) {
	datacomment.Stats.TotalReactions -= 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Stats.Like -= 1
	case sharedEnums.ReactionLove:
		datacomment.Stats.Love -= 1
	case sharedEnums.ReactionHaha:
		datacomment.Stats.Haha -= 1
	case sharedEnums.ReactionSad:
		datacomment.Stats.Sad -= 1
	case sharedEnums.ReactionAngry:
		datacomment.Stats.Angry -= 1
	case sharedEnums.ReactionWow:
		datacomment.Stats.Wow -= 1
	}
}
