package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
)

type IUserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(userID uuid.UUID) error
	FindByKeycloakID(keycloakID string) (*entity.User, error)
}
