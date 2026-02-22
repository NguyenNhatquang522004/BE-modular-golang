package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -----------------------------------------------------------------------------
// REQ MAPPER
// -----------------------------------------------------------------------------

// ToEntityPageRole chuyển từ PageRoleReq sang Entity (Thường dùng cho Create)
func ToEntityPageRole(r *req.PageRoleReq) *entity.PageRole {
	if r == nil {
		return nil
	}

	// Xử lý ID
	objID := primitive.NewObjectID()
	if r.ID != "" {
		if parsedID, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			objID = parsedID
		}
	}

	// Xử lý PageID
	pageID := primitive.NilObjectID
	if r.PageID != "" {
		if parsedPageID, err := primitive.ObjectIDFromHex(r.PageID); err == nil {
			pageID = parsedPageID
		}
	}

	e := &entity.PageRole{
		ID:                objID,
		PageID:            pageID,
		UserID:            r.UserID,
		Role:              r.Role,
		CustomPermissions: copyStringSlice(r.CustomPermissions),
		CreatedAt:         time.Now(), // Mặc định gán thời gian hiện tại khi tạo mới
		AssignedBy:        r.AssignedBy,
	}

	// Nếu request truyền lên CreatedAt cụ thể (ví dụ import data cũ), ta lấy giá trị đó
	if !r.CreatedAt.IsZero() {
		e.CreatedAt = r.CreatedAt
	}

	return e
}

// UpdateToEntityPageRole cập nhật các trường từ PageRoleReq vào Entity có sẵn (Thường dùng cho Update)
func UpdateToEntityPageRole(r *req.PageRoleReq, e *entity.PageRole) {
	if r == nil || e == nil {
		return
	}

	// Cập nhật ID (nếu có sự thay đổi, dù hiếm khi update ID)
	if r.ID != "" {
		if parsedID, err := primitive.ObjectIDFromHex(r.ID); err == nil {
			e.ID = parsedID
		}
	}

	// Cập nhật PageID
	if r.PageID != "" {
		if parsedPageID, err := primitive.ObjectIDFromHex(r.PageID); err == nil {
			e.PageID = parsedPageID
		}
	}

	e.UserID = r.UserID
	e.Role = r.Role
	e.CustomPermissions = copyStringSlice(r.CustomPermissions)
	e.AssignedBy = r.AssignedBy

	// Lưu ý: Thông thường CreatedAt không nên update, nhưng ánh xạ đủ 100% theo yêu cầu
	if !r.CreatedAt.IsZero() {
		e.CreatedAt = r.CreatedAt
	}
}

// -----------------------------------------------------------------------------
// RES MAPPER
// -----------------------------------------------------------------------------

// ToResponse chuyển từ Entity sang PageRoleRes
func ToResponse(e *entity.PageRole) *res.PageRoleRes {
	if e == nil {
		return nil
	}

	return &res.PageRoleRes{
		ID:                e.ID.Hex(),
		PageID:            e.PageID.Hex(),
		UserID:            e.UserID,
		Role:              e.Role,
		CustomPermissions: copyStringSlice(e.CustomPermissions),
		CreatedAt:         e.CreatedAt,
		AssignedBy:        e.AssignedBy,
	}
}

// ReqToResponsePageRole chuyển trực tiếp từ Req sang Res (Theo đúng yêu cầu của bạn)
func ReqToResponsePageRole(r *req.PageRoleReq) *res.PageRoleRes {
	if r == nil {
		return nil
	}

	return &res.PageRoleRes{
		ID:                r.ID,
		PageID:            r.PageID,
		UserID:            r.UserID,
		Role:              r.Role,
		CustomPermissions: copyStringSlice(r.CustomPermissions),
		CreatedAt:         r.CreatedAt,
		AssignedBy:        r.AssignedBy,
	}
}

// -----------------------------------------------------------------------------
// HELPER FUNC
// -----------------------------------------------------------------------------

// copyStringSlice giúp deep copy array, tránh reference memory ngoài ý muốn giữa các layer
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
