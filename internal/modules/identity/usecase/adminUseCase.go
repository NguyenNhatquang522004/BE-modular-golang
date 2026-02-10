package usecase

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryKeyCloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
)

type AdminUseCase struct {
	userRepo        IRepositoryPostgres.IUserRepository
	keycloakClient  IRepositoryKeyCloak.IKeycloakRepository
	userUseCase     IUserService
	UserRoleUseCase IUserRoleUseCase
}

func NewAdminUseCase(userRepo IRepositoryPostgres.IUserRepository, keycloakClient IRepositoryKeyCloak.IKeycloakRepository, userRoleUseCase IUserRoleUseCase, userUseCase IUserService) *AdminUseCase {
	return &AdminUseCase{
		userRepo:        userRepo,
		keycloakClient:  keycloakClient,
		userUseCase:     userUseCase,
		UserRoleUseCase: userRoleUseCase,
	}
}
func (a *AdminUseCase) LoadAllUsers(Cursor string, Limit int) (*response.Response, error) {
	data, err := a.userUseCase.PanigationUsers(Cursor, Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error load all users"), response.WithStatus("500")), err
	}
	return data, nil
	// Implement logic to get all users
}
func (a *AdminUseCase) GetUserByID(userID string) (*response.Response, error) {
	data, err := a.userUseCase.GetUserByID(userID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error get user by id"), response.WithStatus("500")), err
	}
	return data, nil
}
func (a *AdminUseCase) UpdateUser(req *req.UpdateUserReq) (*response.Response, error) {
	// Implement logic to update user information
	data, err := a.userUseCase.UpdateUser(req)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error update user"), response.WithStatus("500")), err
	}
	return data, nil
}
func (a *AdminUseCase) DeleteUser(userID string) (*response.Response, error) {
	ctx := context.Background()
	data, err := a.userUseCase.DeleteUser(userID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error delete user"), response.WithStatus("500")), err
	}
	err = a.keycloakClient.DeleteUser(ctx, userID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error delete user in keycloak"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(data),
		response.WithMessage("delete user successful"), response.WithStatus("200")), nil
}
func (a *AdminUseCase) CreateUser(rolename []enum.RoleType, roleID []string, req *req.CreateUserReq) (*response.Response, error) {
	data, err := a.userUseCase.CreateUser(req)

	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error create user"), response.WithStatus("500")), err
	}
	user := &gocloak.User{
		ID:       gocloak.StringP(data.Data.(map[string]any)["user"].(*entity.User).ID.String()),
		Username: gocloak.StringP(data.Data.(map[string]any)["user"].(*entity.User).Username),
		Email:    gocloak.StringP(data.Data.(map[string]any)["user"].(*entity.User).Email),
		Enabled:  gocloak.BoolP(true),
	}
	datakey, err := a.keycloakClient.CreateUser(context.Background(), user, data.Data.(map[string]any)["user"].(*entity.User).Password)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error create user in keycloak"), response.WithStatus("500")), err
	}
	err = a.keycloakClient.UpdateListRealmRole(context.Background(), datakey, rolename)

	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in keycloak"), response.WithStatus("500")), err
	}
	dataRole, err := a.UserRoleUseCase.UpdateRoleOfUser(data.Data.(map[string]any)["user"].(*entity.User).ID.String(), roleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in postgres"), response.WithStatus("500")), err
	}
	if dataRole.Status == "error" {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in postgres"), response.WithStatus("500")), nil
	}
	return data, nil
}
func (a *AdminUseCase) AssignRoleToUser(userID string, rolename []enum.RoleType, roleID []string) (*response.Response, error) {
	ctx := context.Background()
	err := a.keycloakClient.UpdateListRealmRole(ctx, userID, rolename)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in keycloak"), response.WithStatus("500")), err
	}
	data, err := a.UserRoleUseCase.UpdateRoleOfUser(userID, roleID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in postgres"), response.WithStatus("500")), err
	}
	if data.Status == "error" {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error assign role to user in postgres"), response.WithStatus("500")), nil
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("assign role to user successful"), response.WithStatus("200")), nil
}
