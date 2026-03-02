package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type UpdateGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewUpdateGroup() *UpdateGroup {
	return &UpdateGroup{}
}

func (u *UpdateGroup) Execute(ctx context.Context, req *req.UpdateGroupRequest) (*res.FailedGroup, error) {
	datagroup, err := u.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("failed to get group member"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group not found"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != req.GroupID {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group ID mismatch"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	datamember, err := u.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserID, req.GroupID)
	if err != nil {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("failed to get group member"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datamember.Role != sharedEnums.RoleTypeModerator || datamember.Role != sharedEnums.RoleTypeAdmin {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("user does not have the required role"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	mapper.UpdateToEntityGroup(req.UpdateGroupReq, datagroup)
	err = u.groupRepo.UpdateGroup(ctx, datagroup)
	if err != nil {
		return &res.FailedGroup{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("failed to update group"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedGroup{
		GroupID:      req.GroupID,
		UserID:       req.UserID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
