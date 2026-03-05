package IRepositoryPostgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserSessionRepository interface {
	CreateSession(ctx context.Context, session *entity.UserSession) error
	GetAllUserSessions(ctx context.Context, userID string) ([]*entity.UserSession, error)
	GetUserSessionPast(ctx context.Context, userID string) (*entity.UserSession, error)
}
