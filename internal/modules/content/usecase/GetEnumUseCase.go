package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	res "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
)

type GetEnumUseCase struct{}

func NewGetEnumUseCase() *GetEnumUseCase {
	return &GetEnumUseCase{}
}

func (uc *GetEnumUseCase) Execute(ctx context.Context) (*response.Response, error) {
	data := map[string][]*res.EnumReponse{
		"MediaType":        toEnumReponse(toAny(sharedEnums.MediaTypeValues())),
		"PostType":         toEnumReponse(toAny(sharedEnums.PostTypeValues())),
		"ContextType":      toEnumReponse(toAny(sharedEnums.ContextTypeValues())),
		"PrivacyScope":     toEnumReponse(toAny(sharedEnums.PrivacyScopeValues())),
		"PostStatus":       toEnumReponse(toAny(sharedEnums.PostStatusValues())),
		"ActivityType":     toEnumReponse(toAny(sharedEnums.ActivityTypeValues())),
		"PublisherRole":    toEnumReponse(toAny(sharedEnums.RoleTypeValues())),
		"TargetCollection": toEnumReponse(toAny(sharedEnums.TargetCollectionValues())),
	}
	return response.NewResponse(response.WithData(data)), nil
}

// enumValuer is satisfied by all enumer-generated types (they all have String() and an int underlying type)
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
