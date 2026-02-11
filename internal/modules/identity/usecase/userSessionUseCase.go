package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
)

type UserSessionUseCase struct {
	userSessionRepo IRepositoryPostgres.IUserSessionRepository
}

func NewUserSessionUseCase(userSessionRepo IRepositoryPostgres.IUserSessionRepository) *UserSessionUseCase {
	return &UserSessionUseCase{
		userSessionRepo: userSessionRepo,
	}
}
func (u *UserSessionUseCase) CreateSessionLogin(session *req.UserSessionReq) (*response.Response, error) {
	entity := session.ToEntity()
	if entity == nil {
		return response.NewResponse(response.WithData(""), response.WithMessage("invalid session data"), response.WithStatus("400")), nil
	}
	return u.userSessionRepo.CreateSession(entity)
}
func (u *UserSessionUseCase) GetAllUserSessions(req *req.UserSessionReq) (*response.Response, error) {
	entity := req.ToEntity()
	return u.userSessionRepo.GetAllUserSessions(entity.UserID.String())
}
func (u *UserSessionUseCase) GetUserSessionPast(req *req.UserSessionReq) (*response.Response, error) {
	entity := req.ToEntity()
	return u.userSessionRepo.GetUserSessionPast(entity.UserID.String())
}
