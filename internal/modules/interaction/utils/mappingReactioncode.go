package utils

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

func MappingCreateReactioncode(datacomment *entity.Comment, react *sharedEnums.ReactionCode) {
	datacomment.Reactions.Total += 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Reactions.Like += 1
	case sharedEnums.ReactionLove:
		datacomment.Reactions.Love += 1
	case sharedEnums.ReactionHaha:
		datacomment.Reactions.Haha += 1
	case sharedEnums.ReactionSad:
		datacomment.Reactions.Sad += 1
	case sharedEnums.ReactionAngry:
		datacomment.Reactions.Angry += 1
	case sharedEnums.ReactionWow:
		datacomment.Reactions.Wow += 1
	}
}

func MappingDeleteReactioncode(datacomment *entity.Comment, react *sharedEnums.ReactionCode) {
	datacomment.Reactions.Total -= 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Reactions.Like -= 1
	case sharedEnums.ReactionLove:
		datacomment.Reactions.Love -= 1
	case sharedEnums.ReactionHaha:
		datacomment.Reactions.Haha -= 1
	case sharedEnums.ReactionSad:
		datacomment.Reactions.Sad -= 1
	case sharedEnums.ReactionAngry:
		datacomment.Reactions.Angry -= 1
	case sharedEnums.ReactionWow:
		datacomment.Reactions.Wow -= 1
	}
}
