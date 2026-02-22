package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionPageRoles = "PageRoles"
)

type PageRole struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LIÊN KẾT (RELATIONSHIPS)
	// Page nào? (Mongo ID)
	// Index: Compound Unique { page_id: 1, user_id: 1 } -> 1 User chỉ có 1 vai trò trên 1 Page
	PageID primitive.ObjectID `bson:"page_id" json:"page_id"`

	// User nào? (Postgres UUID -> String)
	// Index: { user_id: 1 } -> Tìm "Các Page tôi đang quản lý"
	UserID string `bson:"user_id" json:"user_id"`

	// 2. PHÂN QUYỀN (RBAC + ABAC)
	// Vai trò chính (Role-Based)
	Role enum.PageRole `bson:"role" json:"role"`

	// Quyền tùy chỉnh thêm (Attribute-Based)
	// VD: ["manage_jobs", "manage_events", "publish_stories"]
	// Dùng để mở rộng quyền cho các role thấp hơn hoặc giới hạn quyền.
	CustomPermissions []string `bson:"custom_permissions,omitempty" json:"custom_permissions,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`

	// Người cấp quyền này (Postgres UUID -> String)
	AssignedBy string `bson:"assigned_by" json:"assigned_by"`
}

func (PageRole) CollectionName() string {
	return CollectionPageRoles
}
