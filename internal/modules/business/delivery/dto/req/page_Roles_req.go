package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// PageRoleReq ánh xạ 100% các trường từ Entity
type PageRoleReq struct {
	ID                string               `json:"id,omitempty"` // Thường để rỗng khi Create, có giá trị khi Update
	PageID            string               `json:"page_id" validate:"required"`
	UserID            string               `json:"user_id" validate:"required"`
	Role              sharedEnums.RoleType `json:"role" validate:"required"`
	CustomPermissions []string             `json:"custom_permissions,omitempty"`
	CreatedAt         time.Time            `json:"created_at,omitempty"`
	AssignedBy        string               `json:"assigned_by" validate:"required"`
}
