package postgres

import (
	"context"
	"fmt"

	irepositoryshare "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	// repository fields
	DB        *gorm.DB
	redisRepo irepositoryshare.IRedis
}

func NewUserRepository(db *gorm.DB, redisRepo irepositoryshare.IRedis) *UserRepository {
	return &UserRepository{
		DB:        db,
		redisRepo: redisRepo,
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
	var user = &entity.User{}
	err := r.DB.Where(&entity.User{Email: email}).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	var user = &entity.User{}
	err := r.DB.Where(&entity.User{ID: uuid.MustParse(userID)}).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
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
	var user = &entity.User{}
	err := r.DB.Where(&entity.User{KeycloakID: keycloakID}).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Panigation(Cursor string, Limit int) ([]*entity.User, string, bool, int, error) {
	var users = []*entity.User{}
	users, cursor, hasNext, limit, err := r.GetPage1DataFromCache(context.Background())
	if err != nil {
		return nil, "", false, 0, err // Lỗi hệ thống Redis
	}
	if users != nil {
		// CACHE HIT: Trả về dữ liệu luôn
		return users, cursor, hasNext, limit, nil
	}
	querylimit := Limit + 1
	// 1. Khởi tạo Query cơ bản
	// Quan trọng: Phải sort cố định để cursor hoạt động đúng
	query := r.DB.Model(&entity.User{}).Where(&entity.User{}).Order("created_at DESC, id DESC").Limit(querylimit)

	if Cursor != "" {
		// 2. Lấy thông tin của bản ghi tại
		time, cursorUUID, err := utils.DecodeCursor(Cursor)
		if err != nil {
			return nil, "", false, 0, err
		}
		// 3. Thêm điều kiện cho Query
		query = query.Where("(created_at < ?) OR (created_at = ? AND id > ?)", time, time, cursorUUID)
		err = query.Find(users).Error
		if err != nil {
			return nil, "", false, 0, err
		}
		hasNext := false
		if len(users) == querylimit {
			hasNext = true
			users = users[:Limit]
		}
		lastUser := users[len(users)-1]
		encodecursor := utils.EncodeCursor(lastUser.CreatedAt, lastUser.ID)
		return users, encodecursor, hasNext, Limit, nil

	}
	// Check cache
	err = query.Find(users).Error
	if err != nil {
		return nil, "", false, 0, err
	}

	hasNext = false
	if len(users) == querylimit {
		hasNext = true
		users = users[:Limit]
	}
	lastUser := users[len(users)-1]
	encodecursor := utils.EncodeCursor(lastUser.CreatedAt, lastUser.ID)
	err = r.CachePage1Data(context.Background(), users, encodecursor, hasNext, Limit)
	if err != nil {
		return nil, "", false, 0, err
	}

	return users, encodecursor, hasNext, Limit, nil
}
func (r *UserRepository) CachePage1Data(ctx context.Context, users interface{}, encodeCursor string, hasNext bool, limit int) error {
	// Danh sách các key và value cần lưu
	// Dùng map hoặc slice struct để dễ quản lý nếu danh sách dài
	items := map[string]interface{}{
		"page1_user":              users,
		"page1_user_encodecursor": encodeCursor,
		"page1_user_hasNext":      hasNext,
		"page1_user_Limit":        limit,
	}

	// Duyệt qua từng phần tử và lưu vào Redis
	for key, value := range items {
		// Set với thời gian hết hạn là 0 (hoặc thay đổi tùy business logic)
		_, err := r.redisRepo.Set(ctx, key, value, 0)
		if err != nil {
			// Nếu lỗi 1 cái thì return lỗi ngay (hoặc log lại tùy bạn)
			return fmt.Errorf("failed to cache key %s: %w", key, err)
		}
	}

	return nil
}

func (r *UserRepository) GetPage1DataFromCache(ctx context.Context) ([]*entity.User, string, bool, int, error) {
	// 1. Thử lấy Key quan trọng nhất (Users)
	userDataWrapper, err := r.redisRepo.Get(ctx, "page1_user")
	if err != nil {
		return nil, "", false, 0, err // Lỗi kết nối Redis
	}
	if userDataWrapper.Data == nil {
		return nil, "", false, 0, nil // Cache miss (Không có dữ liệu, không phải lỗi)
	}

	// 2. Nếu có Users, lấy tiếp các thông tin phụ
	// Lưu ý: Nên xử lý trường hợp một trong các key này bị thiếu (dù hiếm)
	cursorWrapper, _ := r.redisRepo.Get(ctx, "page1_user_encodecursor")
	hasNextWrapper, _ := r.redisRepo.Get(ctx, "page1_user_hasNext")
	limitWrapper, _ := r.redisRepo.Get(ctx, "page1_user_Limit")

	// 3. Ép kiểu an toàn (Type Assertion)
	// Dùng cú pháp: value, ok := data.(Type) để tránh panic
	users, ok1 := userDataWrapper.Data.([]*entity.User)
	cursor, ok2 := cursorWrapper.Data.(string)
	hasNext, ok3 := hasNextWrapper.Data.(bool)
	limit, ok4 := limitWrapper.Data.(int)

	// Nếu bất kỳ dữ liệu nào bị sai kiểu hoặc nil, coi như cache hỏng
	if !ok1 || !ok2 || !ok3 || !ok4 {
		// Log warning ở đây nếu cần
		return nil, "", false, 0, nil // Trả về như cache miss để load lại từ DB
	}

	return users, cursor, hasNext, limit, nil
}
