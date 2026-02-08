package IRepositoryPostgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type IUserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(userID string) error
	FindByKeycloakID(keycloakID string) (*entity.User, error)
}
