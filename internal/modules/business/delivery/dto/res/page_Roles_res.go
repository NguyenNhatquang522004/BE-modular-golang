package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// PageRoleRes trả về dữ liệu cho Client (Các trường ObjectID từ Mongo đều được map sang chuỗi String)
type PageRoleRes struct {
	ID                string               `json:"id"`
	PageID            string               `json:"page_id"`
	UserID            string               `json:"user_id"`
	Role              sharedEnums.RoleType `json:"role"`
	CustomPermissions []string             `json:"custom_permissions,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	AssignedBy        string               `json:"assigned_by"`
}
