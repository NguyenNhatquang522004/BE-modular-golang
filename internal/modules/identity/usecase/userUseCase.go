package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/utils"
)

type UserUseCase struct {
	userRepo  IRepositoryPostgres.IUserRepository
	redisRepo IRepositoryShare.IRedis
}

func NewUserUseCase(userRepo IRepositoryPostgres.IUserRepository, redisRepo IRepositoryShare.IRedis) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		redisRepo: redisRepo,
	}
}
func (u *UserUseCase) GetUserByID(userID string) (*response.Response, error) {
	var user *entity.User
	user, _ = u.userRepo.GetUserByID(userID)
	if user == nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user not found"), response.WithStatus("404")), nil
	}
	return response.NewResponse(response.WithData(map[string]any{
		"user": user,
	}), response.WithMessage("success"), response.WithStatus("200")), nil
}

func (u *UserUseCase) GetUserByEmail(email string) (*response.Response, error) {
	var user *entity.User
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user not found"), response.WithStatus("404")), nil
	}
	return response.NewResponse(response.WithData(map[string]any{
		"user": user,
	}), response.WithMessage("success"), response.WithStatus("200")), nil

}

func (u *UserUseCase) PanigationUsers(Cursor string, Limit int) (*response.Response, error) {
	users, nextCursor, err := u.userRepo.Panigation(Cursor, Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error panigation users"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(&res.PanigationUsersRes{
		Users:      users,
		NextCursor: nextCursor,
	}), response.WithMessage("success"), response.WithStatus("200")), nil
}

func (u *UserUseCase) DeleteUser(userID string) (*response.Response, error) {
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error delete user"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("success delete user"), response.WithStatus("200")), nil

}

func (u *UserUseCase) UpdateUser(req *req.UpdateUserReq) (*response.Response, error) {
	user := req.ToEntity()
	err := u.userRepo.UpdateUser(user)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error update user"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(map[string]any{
		"user": user,
	}), response.WithMessage("success"), response.WithStatus("200")), nil
}

func (u *UserUseCase) CreateUser(req *req.CreateUserReq) (*response.Response, error) {
	user := req.ToEntity()
	passwordhash, checkPasswordErr := utils.HashPassword(req.Password)
	if checkPasswordErr != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("error hash password"), response.WithStatus("500")), checkPasswordErr
	}
	user.Password = passwordhash
	data, err := u.userRepo.CreateUser(user)
	if err != nil {
		return response.NewResponse(response.WithData(data),
			response.WithMessage("error create user"), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(data), response.WithMessage("success"), response.WithStatus("200")), nil
}
