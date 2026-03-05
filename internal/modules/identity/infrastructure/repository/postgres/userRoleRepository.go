package postgres

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}

func (u *UserRoleRepository) CreateRole(ctx context.Context, role sharedEnums.RoleType, description string) error {
	return u.db.WithContext(ctx).Create(&entity.UserRole{
		Role:        role,
		Description: description,
	}).Error
}

func (u *UserRoleRepository) DeleteRole(ctx context.Context, roleID string) error {
	role, err := u.FindRoleWithID(ctx, roleID)
	if err != nil {
		return err
	}
	return u.db.WithContext(ctx).Delete(role).Error
}

func (u *UserRoleRepository) GetAllUserRoles(ctx context.Context, RoleID string) ([]*entity.UserRole, error) {
	var roles []*entity.UserRole
	err := u.db.WithContext(ctx).
		Joins("JOIN user_roles on user_roles.user_id = users.id").
		Joins("JOIN roles on roles.id = user_roles.role_id").
		Where("roles.id = ?", RoleID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *UserRoleRepository) FindUserwithRole(ctx context.Context, userID string, role sharedEnums.RoleType) (*entity.UserRole, error) {
	var userRole = &entity.UserRole{}
	err := u.db.WithContext(ctx).
		Joins("JOIN user_user_roles on user_user_roles.user_role_id = user_roles.id").
		Where("user_user_roles.user_id = ? AND user_roles.role = ?", userID, role).
		First(userRole).Error
	if err != nil {
		return nil, err
	}
	return userRole, nil
}

func (u *UserRoleRepository) FindRoleWithID(ctx context.Context, roleID string) (*entity.UserRole, error) {
	var role = &entity.UserRole{}
	err := u.db.WithContext(ctx).Where("id = ?", roleID).First(role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (u *UserRoleRepository) FindRoleWithName(ctx context.Context, role sharedEnums.RoleType) (*entity.UserRole, error) {
	var userRole = &entity.UserRole{}
	err := u.db.WithContext(ctx).Where("role = ?", role).First(userRole).Error
	if err != nil {
		return nil, err
	}
	return userRole, nil
}

func (u *UserRoleRepository) UpdateRoleDescription(ctx context.Context, roleID string, description string) error {
	role, err := u.FindRoleWithID(ctx, roleID)
	if err != nil {
		return err
	}
	role.Description = description
	return u.db.WithContext(ctx).Save(role).Error
}

func (u *UserRoleRepository) CreateRoleUser(ctx context.Context, userID string, roleID string) error {
	return u.db.WithContext(ctx).
		Model(&entity.User{ID: uuid.MustParse(userID)}).
		Association("Roles").
		Append(&entity.UserRole{ID: uuid.MustParse(roleID)})
}

func (u *UserRoleRepository) GetAllRoles(ctx context.Context) error {
	var roles []*entity.UserRole
	return u.db.WithContext(ctx).Find(&roles).Error
}

func (u *UserRoleRepository) DeleteRoleFromUser(ctx context.Context, userID string, roleID string) error {
	return u.db.WithContext(ctx).
		Model(&entity.User{ID: uuid.MustParse(userID)}).
		Association("Roles").
		Delete(&entity.UserRole{ID: uuid.MustParse(roleID)})
}

func (u *UserRoleRepository) UpdateRoleOfUser(ctx context.Context, userID string, roleIDs []string) error {
	var roles []*entity.UserRole
	for _, id := range roleIDs {
		roles = append(roles, &entity.UserRole{ID: uuid.MustParse(id)})
	}
	return u.db.WithContext(ctx).
		Model(&entity.User{ID: uuid.MustParse(userID)}).
		Association("Roles").
		Replace(roles)
}
