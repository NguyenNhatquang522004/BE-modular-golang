package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	interactionEnum "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
)

type GetEnumUsecase struct {
}

func NewGetEnumUsecase() *GetEnumUsecase {
	return &GetEnumUsecase{}
}

func (u *GetEnumUsecase) Execute(ctx context.Context) (*response.Response, error) {
	data := map[string][]*res.EnumReponse{
		// interaction/enum - Comment
		"CommentStatus": toEnumReponse(toAny(interactionEnum.CommentStatusValues())),

		// sharedEnums - dùng trong entity_reactions, reaction_history
		"ReactionTarget": toEnumReponse(toAny(sharedEnums.ReactionTargetValues())),
		"ReactionCode":   toEnumReponse(toAny(sharedEnums.ReactionCodeValues())),

		// sharedEnums - dùng trong comment_Edit_Logs
		"TargetCollection": toEnumReponse(toAny(sharedEnums.TargetCollectionValues())),

		// sharedEnums - dùng trong user_Saved_Items
		"SavedTargetType": toEnumReponse(toAny(sharedEnums.SavedTargetTypeValues())),
	}
	return response.NewResponse(response.WithData(data)), nil
}

type enumValuer interface {
	String() string
	int() int
}

// enumItem wraps any enum value with its int cast and string label
type enumItem struct {
	val int
	str string
}

func toAny[T ~int](values []T) []enumItem {
	items := make([]enumItem, 0, len(values))
	for _, v := range values {
		type stringer interface{ String() string }
		var s string
		if sv, ok := any(v).(stringer); ok {
			s = sv.String()
		}
		items = append(items, enumItem{val: int(v), str: s})
	}
	return items
}

func toEnumReponse(items []enumItem) []*res.EnumReponse {
	result := make([]*res.EnumReponse, 0, len(items))
	for _, item := range items {
		result = append(result, &res.EnumReponse{
			Value: item.val,
			Label: item.str,
		})
	}
	return result
}
