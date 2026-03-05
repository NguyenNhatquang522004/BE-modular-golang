package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserRoleRepository interface {
	CreateRole(ctx context.Context, role sharedEnums.RoleType, description string) error
	DeleteRole(ctx context.Context, roleID string) error
	GetAllUserRoles(ctx context.Context, RoleID string) ([]*entity.UserRole, error)
	FindUserwithRole(ctx context.Context, userID string, role sharedEnums.RoleType) (*entity.UserRole, error)
	FindRoleWithID(ctx context.Context, roleID string) (*entity.UserRole, error)
	FindRoleWithName(ctx context.Context, role sharedEnums.RoleType) (*entity.UserRole, error)
	UpdateRoleDescription(ctx context.Context, roleID string, description string) error
	CreateRoleUser(ctx context.Context, userID string, roleID string) error
	GetAllRoles(ctx context.Context) error
	DeleteRoleFromUser(ctx context.Context, userID string, roleID string) error
	UpdateRoleOfUser(ctx context.Context, userID string, roleIDs []string) error
}
