package usecase

import (
	"context"
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
)

type UserRoleUseCase struct {
	userRoleRepo IRepositoryPostgres.IUserRoleRepository
	userRepo     IRepositoryPostgres.IUserRepository
}

func NewUserRoleUseCase(userRoleRepo IRepositoryPostgres.IUserRoleRepository, userRepo IRepositoryPostgres.IUserRepository) *UserRoleUseCase {
	return &UserRoleUseCase{
		userRoleRepo: userRoleRepo,
		userRepo:     userRepo,
	}
}

func (u *UserRoleUseCase) convertEnum(roleName string) (sharedEnums.RoleType, error) {
	normalizedRole := strings.ToLower(strings.TrimSpace(roleName))
	return sharedEnums.RoleTypeString(normalizedRole)
}

func (u *UserRoleUseCase) AssignRoleToUserWithName(ctx context.Context, userID string, roleName string) (*response.Response, error) {
	roleType, err := u.convertEnum(roleName)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("invalid role name"), response.WithStatus("400")), err
	}
	role, err := u.userRoleRepo.FindRoleWithName(ctx, roleType)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("role not found"), response.WithStatus("404")), err
	}
	err = u.userRoleRepo.CreateRoleUser(ctx, userID, role.ID.String())
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("assign fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) AssignRoleToUserWithID(ctx context.Context, userID string, roleID string) (*response.Response, error) {
	err := u.userRoleRepo.CreateRoleUser(ctx, userID, roleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("assign fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) UpdateRoleOfUser(ctx context.Context, userID string, roleIDs []string) (*response.Response, error) {
	err := u.userRoleRepo.UpdateRoleOfUser(ctx, userID, roleIDs)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("update fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role updated for user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) GetUserRoles(ctx context.Context, userID string, roleName string) (*response.Response, error) {
	roleType, err := u.convertEnum(roleName)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("invalid role name"), response.WithStatus("400")), err
	}
	userRole, err := u.userRoleRepo.FindUserwithRole(ctx, userID, roleType)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("role not found"), response.WithStatus("404")), err
	}
	return response.NewResponse(response.WithData(userRole), response.WithMessage(""), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) UpdateRoleDescription(ctx context.Context, roleID string, description string) (*response.Response, error) {
	err := u.userRoleRepo.UpdateRoleDescription(ctx, roleID, description)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("update fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role description updated successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) CreateRole(ctx context.Context, role sharedEnums.RoleType, description string) (*response.Response, error) {
	err := u.userRoleRepo.CreateRole(ctx, role, description)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("create fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role created successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) DeleteRole(ctx context.Context, roleID string) (*response.Response, error) {
	err := u.userRoleRepo.DeleteRole(ctx, roleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("delete fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role deleted successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) GetAllUserRoles(ctx context.Context, RoleID string) (*response.Response, error) {
	roles, err := u.userRoleRepo.GetAllUserRoles(ctx, RoleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("get all user roles fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(roles), response.WithMessage(""), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) GetAllRoles(ctx context.Context) (*response.Response, error) {
	err := u.userRoleRepo.GetAllRoles(ctx)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("get all roles fail"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage(""), response.WithStatus("200")), nil
}
