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
	err := r.DB.Where("id = ?", userID).Delete(&entity.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByKeycloakID(keycloakID string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where(&entity.User{KeycloakID: keycloakID}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Panigation(Cursor string, Limit int) ([]*entity.User, string, error) {
	var users []*entity.User
	// 1. Khởi tạo Query cơ bản
	// Quan trọng: Phải sort cố định để cursor hoạt động đúng
	query := r.DB.Model(&entity.User{}).Where(&entity.User{}).Order("created_at DESC, id DESC").Limit(Limit)

	if Cursor != "" {
		// 2. Lấy thông tin của bản ghi tại cursor
		var cursorUser entity.User
		err := r.DB.Where("id = ?", Cursor).First(&cursorUser).Error
		if err != nil {
			return nil, "", err
		}
		// 3. Thêm điều kiện để lấy các bản ghi sau cursor
		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", cursorUser.CreatedAt, cursorUser.CreatedAt, cursorUser.ID)
	}
	// 4. Thực hiện truy vấn
	err := query.Find(&users).Error
	if err != nil {
		return nil, "", err
	}
	lastUser := ""
	// 5. Xác định next cursor
	if len(users) == Limit {
		lastUser = users[len(users)-1].ID.String()
	}
	return users, lastUser, nil
}
