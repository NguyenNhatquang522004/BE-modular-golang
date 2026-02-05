package IRepositoryPostgres

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"

type IUserSessionRepository interface {
	GetActiveSessions(userID string) (*response.Response, error)
	RevokeSession(sessionID string) (*response.Response, error)
}
