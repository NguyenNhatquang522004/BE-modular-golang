package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type IUserRoleRepository interface {
	CreateRole(ctx context.Context, role sharedEnums.RoleType, description string) (*response.Response, error)
	DeleteRole(ctx context.Context, roleID string) (*response.Response, error)
	GetAllUserRoles(ctx context.Context, RoleID string) (*response.Response, error)
	FindUserwithRole(ctx context.Context, userID string, role sharedEnums.RoleType) (*response.Response, error)
	FindRoleWithID(ctx context.Context, roleID string) (*response.Response, error)
	FindRoleWithName(ctx context.Context, role sharedEnums.RoleType) (*response.Response, error)
	UpdateRoleDescription(ctx context.Context, roleID string, description string) (*response.Response, error)
	CreateRoleUser(ctx context.Context, userID string, roleID string) (*response.Response, error)
	GetAllRoles(ctx context.Context) (*response.Response, error)
	DeleteRoleFromUser(ctx context.Context, userID string, roleID string) (*response.Response, error)
	UpdateRoleOfUser(ctx context.Context, userID string, roleIDs []string) (*response.Response, error)
}
