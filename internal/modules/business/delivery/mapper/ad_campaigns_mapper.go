package mapper

import (
	"gorm.io/gorm"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/google/uuid"
	// Sửa lại đường dẫn import tương ứng với module của bạn
)

// ==========================================
// MAPPER REQ TO ENTITY
// ==========================================

// ToEntityAdCampaign chuyển đổi từ DTO Request sang Entity để lưu mới vào DB.
func ToEntityAdCampaign(r req.AdCampaignReq) entity.AdCampaign {
	campaign := entity.AdCampaign{
		ID:             r.CampaignID,
		AccountID:      r.AccountID,
		Name:           r.Name,
		Objective:      *r.Objective,
		BuyingType:     *r.BuyingType,
		DailyBudget:    r.DailyBudget,
		LifetimeBudget: r.LifetimeBudget,
		StartTime:      r.StartTime,
		EndTime:        r.EndTime,
		Status:         *r.Status,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	if r.DeletedAt != nil {
		campaign.DeletedAt = gorm.DeletedAt{Time: *r.DeletedAt, Valid: true}
	}

	return campaign
}

// UpdateToEntityAdCampaign ánh xạ dữ liệu từ DTO Request sang Entity đang tồn tại.
// Giả định các enum (CampaignObjective, BuyingType, CampaignStatus) có underlying type là string.
func UpdateToEntityAdCampaign(r req.AdCampaignReq, e *entity.AdCampaign) {
	if r.CampaignID != uuid.Nil {
		e.ID = r.CampaignID
	}
	if r.AccountID != uuid.Nil {
		e.AccountID = r.AccountID
	}
	if r.Name != "" {
		e.Name = r.Name
	}
	if r.Objective != nil {
		e.Objective = *r.Objective
	}
	if r.BuyingType != nil {
		e.BuyingType = *r.BuyingType
	}
	// Với kiểu con trỏ, so sánh thẳng với nil
	if r.DailyBudget != nil {
		e.DailyBudget = r.DailyBudget
	}
	if r.LifetimeBudget != nil {
		e.LifetimeBudget = r.LifetimeBudget
	}
	if !r.StartTime.IsZero() {
		e.StartTime = r.StartTime
	}
	if r.EndTime != nil {
		e.EndTime = r.EndTime
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if !r.CreatedAt.IsZero() {
		e.CreatedAt = r.CreatedAt
	}
	if !r.UpdatedAt.IsZero() {
		e.UpdatedAt = r.UpdatedAt
	}
	if r.DeletedAt != nil {
		e.DeletedAt = gorm.DeletedAt{Time: *r.DeletedAt, Valid: true}
	}
}

// ==========================================
// MAPPER TO RES
// ==========================================

// ToResFromReqAdCampaign chuyển đổi trực tiếp từ Req sang Res theo đúng yêu cầu của bạn.
func ToResFromReqAdCampaign(r req.AdCampaignReq) res.AdCampaignRes {
	return res.AdCampaignRes{
		CampaignID:     r.CampaignID,
		AccountID:      r.AccountID,
		Name:           r.Name,
		Objective:      r.Objective,
		BuyingType:     r.BuyingType,
		DailyBudget:    r.DailyBudget,
		LifetimeBudget: r.LifetimeBudget,
		StartTime:      r.StartTime,
		EndTime:        r.EndTime,
		Status:         r.Status,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
		DeletedAt:      r.DeletedAt,
	}
}

// ToResFromEntityAdCampaign (Tặng kèm thêm hàm Best Practice): Dùng để map từ Entity lấy dưới DB lên trả về Res.
func ToResFromEntityAdCampaign(e entity.AdCampaign) res.AdCampaignRes {
	response := res.AdCampaignRes{
		CampaignID:     e.ID,
		AccountID:      e.AccountID,
		Name:           e.Name,
		Objective:      &e.Objective,
		BuyingType:     &e.BuyingType,
		DailyBudget:    e.DailyBudget,
		LifetimeBudget: e.LifetimeBudget,
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		Status:         &e.Status,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}

	if e.DeletedAt.Valid {
		response.DeletedAt = &e.DeletedAt.Time
	}

	return response
}
