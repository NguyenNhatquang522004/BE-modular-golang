package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{
		db: db,
	}
}

func (r *UserSessionRepository) CreateSession(ctx context.Context, session *entity.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *UserSessionRepository) GetAllUserSessions(ctx context.Context, userID string) ([]*entity.UserSession, error) {
	var sessions []*entity.UserSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *UserSessionRepository) GetUserSessionPast(ctx context.Context, userID string) (*entity.UserSession, error) {
	var session entity.UserSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Offset(1).Limit(1).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
