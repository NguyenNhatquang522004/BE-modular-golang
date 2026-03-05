package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSessionRepository interface {
	CreateSession(ctx context.Context, session *entity.UserSession) (*response.Response, error)
	GetAllUserSessions(ctx context.Context, userID string) (*response.Response, error)
	GetUserSessionPast(ctx context.Context, userID string) (*response.Response, error)
}
