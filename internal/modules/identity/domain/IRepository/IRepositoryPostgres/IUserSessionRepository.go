package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSessionRepository interface {
	CreateSession(session *entity.UserSession) (*response.Response, error)
	GetAllUserSessions(userID string) (*response.Response, error)
	GetUserSessionPast(userID string) (*response.Response, error)
}
