package postgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	// repository fields
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	user.ID = uuid.New()
	err := r.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where(&entity.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where(&entity.User{ID: uuid.MustParse(userID)}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) UpdateUser(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) DeleteUser(userID string) error {
	return r.DB.Where("id = ?", userID).Delete(&entity.User{}).Error
}

func (r *UserRepository) FindByKeycloakID(keycloakID string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where(&entity.User{KeycloakID: keycloakID}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
