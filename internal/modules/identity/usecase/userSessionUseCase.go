package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type UserSessionUseCase struct {
	userSessionRepo IRepositoryPostgres.IUserSessionRepository
}

func NewUserSessionUseCase(userSessionRepo IRepositoryPostgres.IUserSessionRepository) *UserSessionUseCase {
	return &UserSessionUseCase{
		userSessionRepo: userSessionRepo,
	}
}
func (u *UserSessionUseCase) CreateSessionLogin(session *entity.UserSession) (*response.Response, error) {
	return u.userSessionRepo.CreateSession(session)
}
func (u *UserSessionUseCase) GetAllUserSessions(userID string) (*response.Response, error) {
	return u.userSessionRepo.GetAllUserSessions(userID)
}
func (u *UserSessionUseCase) GetUserSessionPast(userID string) (*response.Response, error) {
	return u.userSessionRepo.GetUserSessionPast(userID)
}
