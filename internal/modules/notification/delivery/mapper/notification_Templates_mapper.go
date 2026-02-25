package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToEntityNotificationTemplate: Ánh xạ từ Create Req sang Entity mới
func ToEntityNotificationTemplate(r *req.CreateNotificationTemplateReq) *entity.NotificationTemplate {
	if r == nil {
		return nil
	}

	// Best Practice: Tự động gán thời gian nếu Request không truyền lên
	now := time.Now()
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

	return &entity.NotificationTemplate{
		ID:         primitive.NewObjectID(),
		Type:       r.Type,
		Template:   r.Template,
		IconURL:    r.IconURL,
		ActionLink: r.ActionLink,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// UpdateToEntityNotificationTemplate: Cập nhật từ Update Req vào Entity hiện có
func UpdateToEntityNotificationTemplate(r *req.UpdateNotificationTemplateReq, ent *entity.NotificationTemplate) {
	if r == nil || ent == nil {
		return
	}

	if r.Type != nil {
		ent.Type = *r.Type
	}
	if r.Template != nil {
		ent.Template = r.Template // Ghi đè toàn bộ map hoặc tự merge key tùy logic nghiệp vụ
	}
	if r.IconURL != nil {
		ent.IconURL = *r.IconURL
	}
	if r.ActionLink != nil {
		ent.ActionLink = *r.ActionLink
	}

	// Best Practice: Tự động cập nhật UpdatedAt khi có thao tác update
	if r.UpdatedAt != nil && !r.UpdatedAt.IsZero() {
		ent.UpdatedAt = *r.UpdatedAt
	} else {
		ent.UpdatedAt = time.Now()
	}
}

// ReqToResNotificationTemplate: Chuyển trực tiếp từ Request sang Response (Theo đúng yêu cầu của bạn)
func ReqToResNotificationTemplate(r *req.CreateNotificationTemplateReq) *res.NotificationTemplateRes {
	if r == nil {
		return nil
	}

	return &res.NotificationTemplateRes{
		ID:         primitive.NewObjectID().Hex(), // Mock một ID mới
		Type:       r.Type,
		Template:   r.Template,
		IconURL:    r.IconURL,
		ActionLink: r.ActionLink,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// EntityToResNotificationTemplate: Chuyển từ Entity DB sang Response (Hàm chuẩn thường dùng trong luồng Get/List)
func EntityToResNotificationTemplate(ent *entity.NotificationTemplate) *res.NotificationTemplateRes {
	if ent == nil {
		return nil
	}

	return &res.NotificationTemplateRes{
		ID:         ent.ID.Hex(),
		Type:       ent.Type,
		Template:   ent.Template,
		IconURL:    ent.IconURL,
		ActionLink: ent.ActionLink,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}

// ToNotificationTemplateRes chuyển đổi 1 entity sang response DTO
func ToNotificationTemplateRes(e entity.NotificationTemplate) res.NotificationTemplateRes {
	return res.NotificationTemplateRes{
		ID:         e.ID.Hex(), // Chuyển ObjectID sang string
		Type:       e.Type,
		Template:   e.Template,
		IconURL:    e.IconURL,
		ActionLink: e.ActionLink,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// ToNotificationTemplateResList chuyển đổi một danh sách entity sang danh sách response DTO (rất tiện khi viết hàm GET list)
func ToNotificationTemplateResList(entities []entity.NotificationTemplate) []res.NotificationTemplateRes {
	result := make([]res.NotificationTemplateRes, len(entities))
	for i, e := range entities {
		result[i] = ToNotificationTemplateRes(e)
	}
	return result
}
