package postgres

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
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

func (u *UserRoleRepository) CreateRole(role enum.RoleType, description string) (*response.Response, error) {
	err := u.db.Create(&entity.UserRole{
		Role:        role,
		Description: description,
	}).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(err), response.WithMessage(""), response.WithStatus("200")), nil
}

func (u *UserRoleRepository) DeleteRole(roleID string) (*response.Response, error) {
	// Implementation goes here
	data, err := u.FindRoleWithID(roleID)
	if err != nil {
		return nil, err
	}
	role := data.Data.(*entity.UserRole)
	err = u.db.Delete(&role).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role deleted successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleRepository) GetAllUserRoles(RoleID string) (*response.Response, error) {
	var user = &[]*entity.User{}
	result := u.db.Joins("JOIN user_roles on user_roles.user_id = users.id").
		Joins("JOIN roles on roles.id = user_roles.role_id").Where("roles.id = ?", RoleID).Find(user).Error
	if result != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(result.Error()), response.WithStatus("500")), result
	}

	return response.NewResponse(response.WithData(user), response.WithMessage(""), response.WithStatus("200")), nil

}

func (u *UserRoleRepository) FindUserwithRole(userID string, role enum.RoleType) (*response.Response, error) {
	var user = &entity.User{}
	err := u.db.Joins("JOIN user_roles on user_roles.user_id = users.id").
		Joins("JOIN roles on roles.id = user_roles.role_id").
		Where("users.id = ? AND roles.role = ?", userID, role).
		First(user).Error
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(user), response.WithMessage(""), response.WithStatus("200")), nil
}
func (u *UserRoleRepository) FindRoleWithID(roleID string) (*response.Response, error) {
	var role = &entity.UserRole{}
	err := u.db.Where("id = ?", roleID).First(role).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(role), response.WithMessage(""), response.WithStatus("200")), nil
}
func (u *UserRoleRepository) FindRoleWithName(role enum.RoleType) (*response.Response, error) {
	var userRole = &entity.UserRole{}
	err := u.db.Where("role = ?", role).First(userRole).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(userRole), response.WithMessage(""), response.WithStatus("200")), nil

}

func (u *UserRoleRepository) UpdateRoleDescription(roleID string, description string) (*response.Response, error) {
	role, err := u.FindRoleWithID(roleID)
	if err != nil {
		return nil, err
	}
	userRole := role.Data.(*entity.UserRole)
	userRole.Description = description
	err = u.db.Save(userRole).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(userRole), response.WithMessage(""), response.WithStatus("200")), nil

}
func (u *UserRoleRepository) CreateRoleUser(userID string, roleID string) (*response.Response, error) {
	err := u.db.Model(&entity.User{ID: uuid.MustParse(userID)}).Association("Roles").Append(&entity.UserRole{ID: uuid.MustParse(roleID)})
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("Role assigned to user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleRepository) GetAllRoles() (*response.Response, error) {
	var roles = &[]*res.UserRoleRes{}
	err := u.db.Raw("SELECT ID , role, description FROM user_roles").Scan(roles).Error
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(roles), response.WithMessage(""), response.WithStatus("200")), nil
}

func (u *UserRoleRepository) DeleteRoleFromUser(userID string, roleID string) (*response.Response, error) {
	err := u.db.Model(&entity.User{ID: uuid.MustParse(userID)}).Association("Roles").Delete(&entity.UserRole{ID: uuid.MustParse(roleID)})
	if err != nil {
		return nil, err
	}

	return response.NewResponse(response.WithData(nil), response.WithMessage("Role removed from user successfully"), response.WithStatus("200")), nil
}

func (u *UserRoleRepository) UpdateRoleOfUser(userID string, roleIDs []string) (*response.Response, error) {
	var roles = []*entity.UserRole{}
	for _, id := range roleIDs {
		roles = append(roles, &entity.UserRole{ID: uuid.MustParse(id)})
	}
	err := u.db.Model(&entity.User{ID: uuid.MustParse(userID)}).Association("Roles").Replace(roles)
	if err != nil {
		return nil, err
	}
	return response.NewResponse(response.WithData(nil), response.WithMessage("User roles updated successfully"), response.WithStatus("200")), nil
}
