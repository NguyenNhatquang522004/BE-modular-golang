package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/google/uuid"
	// Sửa lại đường dẫn import tương ứng với module của bạn
)

// ==========================================
// MAPPER REQ TO ENTITY
// ==========================================

// ToEntityAdAccount chuyển đổi toàn bộ từ DTO Request sang Entity (Dùng cho Create).
// Vì DeletedAt ở Entity đã là *time.Time nên ta có thể gán trực tiếp.
func ToEntityAdAccount(r req.AdAccountReq) entity.AdAccount {
	return entity.AdAccount{
		ID:          r.AccountID,
		OwnerUserID: r.OwnerUserID,
		Currency:    r.Currency,
		Timezone:    r.Timezone,
		Balance:     r.Balance,
		CreditLimit: r.CreditLimit,
		Status:      *r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

// UpdateToEntity ánh xạ dữ liệu từ DTO Request sang Entity đang tồn tại (Dùng cho Update).
// Giả định enum.AccountStatus có underlying type là string.
func UpdateToEntityAdAccount(r req.AdAccountReq, e *entity.AdAccount) {
	if r.AccountID != uuid.Nil {
		e.ID = r.AccountID
	}
	if r.OwnerUserID != uuid.Nil {
		e.OwnerUserID = r.OwnerUserID
	}
	if r.Currency != "" {
		e.Currency = r.Currency
	}
	if r.Timezone != "" {
		e.Timezone = r.Timezone
	}

	// Lưu ý Best Practice: Nếu client muốn update Balance/CreditLimit về đúng số 0,
	// điều kiện != 0 sẽ bỏ qua. Để giải quyết triệt để case này trong tương lai,
	// bạn nên cân nhắc đổi 2 trường này trong file req.AdAccountReq thành con trỏ (*float64).
	if r.CreditLimit != 0 {
		e.CreditLimit = r.CreditLimit
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
		e.DeletedAt = r.DeletedAt
	}
}

// ==========================================
// MAPPER TO RES
// ==========================================

// ToResFromReq chuyển đổi trực tiếp từ Req sang Res.
func ToResFromReqAdAccount(r req.AdAccountReq) res.AdAccountRes {
	return res.AdAccountRes{
		AccountID:   r.AccountID,
		OwnerUserID: r.OwnerUserID,
		Currency:    r.Currency,
		Timezone:    r.Timezone,
		Balance:     r.Balance,
		CreditLimit: r.CreditLimit,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

// ToResFromEntity chuyển đổi từ Entity sang Res để trả dữ liệu chuẩn từ DB về cho Client.
func ToResFromEntityAdAccount(e entity.AdAccount) res.AdAccountRes {
	return res.AdAccountRes{
		AccountID:   e.ID,
		OwnerUserID: e.OwnerUserID,
		Currency:    e.Currency,
		Timezone:    e.Timezone,
		Balance:     e.Balance,
		CreditLimit: e.CreditLimit,
		Status:      &e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
	}
}
