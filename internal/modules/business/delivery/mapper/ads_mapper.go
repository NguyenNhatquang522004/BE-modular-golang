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

// ToEntityAds chuyển đổi từ DTO Request sang Entity để lưu vào Database
func ToEntityAds(r req.AdReq) entity.Ad {
	ad := entity.Ad{
		ID:              r.AdID,
		CampaignID:      r.CampaignID,
		TargetPostID:    r.TargetPostID,
		BidAmount:       r.BidAmount,
		Status:          *r.Status,
		RejectionReason: r.RejectionReason,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}

	// Xử lý chuyển đổi pointer time.Time sang gorm.DeletedAt
	if r.DeletedAt != nil {
		ad.DeletedAt = gorm.DeletedAt{Time: *r.DeletedAt, Valid: true}
	}

	return ad
}

// UpdateToEntityAds map dữ liệu từ DTO Request vào một Entity đang có sẵn (dùng cho API Update)
// Chỉ update các field có giá trị hợp lệ để tránh ghi đè dữ liệu rỗng lên DB.
func UpdateToEntityAds(r req.AdReq, e *entity.Ad) {
	if r.AdID != uuid.Nil {
		e.ID = r.AdID
	}
	if r.CampaignID != uuid.Nil {
		e.CampaignID = r.CampaignID
	}
	if r.TargetPostID != "" {
		e.TargetPostID = r.TargetPostID
	}
	if r.BidAmount > 0 {
		e.BidAmount = r.BidAmount
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if r.RejectionReason != nil {
		e.RejectionReason = r.RejectionReason
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
// MAPPER TO RES
// ==========================================

// ToResFromReqAds chuyển đổi trực tiếp từ Req sang Res (Theo đúng yêu cầu của bạn)
func ToResFromReqAds(r req.AdReq) res.AdRes {
	return res.AdRes{
		AdID:            r.AdID,
		CampaignID:      r.CampaignID,
		TargetPostID:    r.TargetPostID,
		BidAmount:       r.BidAmount,
		Status:          r.Status,
		RejectionReason: r.RejectionReason,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		DeletedAt:       r.DeletedAt,
	}
}

// ToResFromEntity chuyển đổi từ Entity sang Res (Best Practice: Thường dùng cái này để trả data từ DB về cho Client)
func ToResFromEntityAds(e entity.Ad) res.AdRes {
	response := res.AdRes{
		AdID:            e.ID,
		CampaignID:      e.CampaignID,
		TargetPostID:    e.TargetPostID,
		BidAmount:       e.BidAmount,
		Status:          &e.Status,
		RejectionReason: e.RejectionReason,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}

	// Xử lý chuyển đổi gorm.DeletedAt sang pointer time.Time
	if e.DeletedAt.Valid {
		response.DeletedAt = &e.DeletedAt.Time
	}

	return response
}
