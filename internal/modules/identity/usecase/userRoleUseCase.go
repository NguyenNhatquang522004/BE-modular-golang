package usecase

import (
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
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

func (u *UserRoleUseCase) AssignRoleToUser(userID string, RoleName string) (*response.Response, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user not found"), response.WithStatus("404")), err
	}
	normalizedRole := strings.ToLower(strings.TrimSpace(RoleName))
	convert, err := enum.RoleTypeString(normalizedRole)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("invalid role name"), response.WithStatus("400")), err
	}
	Role, err := u.userRoleRepo.FindRoleWithName(convert)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("role not found"), response.WithStatus("404")), err
	}
	_, err = u.userRoleRepo.CreateRoleUser(user, Role.Data.(*entity.UserRole))

	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("assign fail "), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleUseCase) RemoveRoleFromUser(userID string, role string) (*response.Response, error) {
	// Implementation here
	return nil, nil
}

func (u *UserRoleUseCase) GetUserRoles(userID string) (*response.Response, error) {
	
	return nil, nil
}
