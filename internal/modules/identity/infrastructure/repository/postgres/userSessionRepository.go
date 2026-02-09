package postgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/req"
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
func (r *UserSessionRepository) CreateSession(session *entity.UserSession) (*response.Response, error) {
	err := r.db.Create(session).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(session), response.WithMessage("success"), response.WithStatus("success")), nil
}

func (r *UserSessionRepository) GetAllUserSessions(userID string) (*response.Response, error) {
	var sessions []req.UserSessionReq
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(sessions), response.WithMessage("success"), response.WithStatus("success")), nil
}
func (r *UserSessionRepository) GetUserSessionPast(userID string) (*response.Response, error) {
	var sessions req.UserSessionReq
	err := r.db.Where("user_id", userID).Order("created_at DESC").Offset(1).Limit(1).First(&sessions).Error
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(sessions), response.WithMessage("success"), response.WithStatus("success")), nil
}
