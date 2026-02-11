package usecase

import (
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
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
func (u *UserRoleUseCase) convertEnum(roleName string) (*response.Response, error) {
	normalizedRole := strings.ToLower(strings.TrimSpace(roleName))
	convert, err := enum.RoleTypeString(normalizedRole)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("invalid role name"), response.WithStatus("400")), err
	}
	return response.NewResponse(response.WithData(convert), response.WithMessage(""), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) AssignRoleToUserWithName(userID string, roleName string) (*response.Response, error) {
	convertResp, err := u.convertEnum(roleName)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("invalid role name"), response.WithStatus("400")), err
	}
	Role, err := u.userRoleRepo.FindRoleWithName(convertResp.Data.(enum.RoleType))
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("role not found"), response.WithStatus("404")), err
	}
	_, err = u.userRoleRepo.CreateRoleUser(userID, Role.Data.(*entity.UserRole).ID.String())

	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("assign fail "), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}
func (u *UserRoleUseCase) AssignRoleToUserWithID(userID string, roleID string) (*response.Response, error) {
	_, err := u.userRoleRepo.CreateRoleUser(userID, roleID)

	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("assign fail "), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) UpdateRoleOfUser(userID string, roleIDs []string) (*response.Response, error) {
	return u.userRoleRepo.UpdateRoleOfUser(userID, roleIDs)
}

func (u *UserRoleUseCase) GetUserRoles(userID string, roleName string) (*response.Response, error) {
	convertResp, err := u.convertEnum(roleName)
	if err != nil {
		return convertResp, err
	}
	return u.userRoleRepo.FindUserwithRole(userID, convertResp.Data.(enum.RoleType))
}

func (u *UserRoleUseCase) UpdateRoleDescription(roleID string, description string) (*response.Response, error) {
	return u.userRoleRepo.UpdateRoleDescription(roleID, description)
}

func (u *UserRoleUseCase) CreateRole(role enum.RoleType, description string) (*response.Response, error) {
	return u.userRoleRepo.CreateRole(role, description)
}

func (u *UserRoleUseCase) DeleteRole(roleID string) (*response.Response, error) {
	return u.userRoleRepo.DeleteRole(roleID)
}

func (u *UserRoleUseCase) GetAllUserRoles(RoleID string) (*response.Response, error) {
	return u.userRoleRepo.GetAllUserRoles(RoleID)
}

func (u *UserRoleUseCase) GetAllRoles() (*response.Response, error) {
	resp, err := u.userRoleRepo.GetAllRoles()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
