package postgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
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
func (r *UserSessionRepository) CreateSession( session *entity.UserSession) (*response.Response, error) {
	err := r.db.Create(session).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}
