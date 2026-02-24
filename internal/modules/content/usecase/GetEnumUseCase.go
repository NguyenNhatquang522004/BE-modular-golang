package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	res "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	enumcontent "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

type GetEnumUseCase struct{}

func NewGetEnumUseCase() *GetEnumUseCase {
	return &GetEnumUseCase{}
}

func (uc *GetEnumUseCase) Execute(ctx context.Context) (*response.Response, error) {
	data := map[string][]*res.EnumReponse{
		"MediaType":        toEnumReponse(toAny(enumcontent.MediaTypeValues())),
		"PostType":         toEnumReponse(toAny(enumcontent.PostTypeValues())),
		"ContextType":      toEnumReponse(toAny(enumcontent.ContextTypeValues())),
		"PrivacyScope":     toEnumReponse(toAny(enumcontent.PrivacyScopeValues())),
		"PostStatus":       toEnumReponse(toAny(enumcontent.PostStatusValues())),
		"ActivityType":     toEnumReponse(toAny(enumcontent.ActivityTypeValues())),
		"PublisherRole":    toEnumReponse(toAny(enumcontent.PublisherRoleValues())),
		"TargetCollection": toEnumReponse(toAny(enumcontent.TargetCollectionValues())),
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
