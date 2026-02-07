package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
)

type IUserRoleRepository interface {
	CreateRole(role enum.RoleType, description string) (*response.Response, error)
	DeleteRole(roleID string) (*response.Response, error)
	GetAllUserRoles(RoleID string) (*response.Response, error)
	FindUserwithRole(userID string, role enum.RoleType) (*response.Response, error)
	FindRoleWithID(roleID string) (*response.Response, error)
	FindRoleWithName(role enum.RoleType) (*response.Response, error)
	UpdateRoleDescription(roleID string, description string) (*response.Response, error)
	CreateRoleUser(user *entity.User, role *entity.UserRole) (*response.Response, error)
	GetAllRoles() (*response.Response, error)
}
