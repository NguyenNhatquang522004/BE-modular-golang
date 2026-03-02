package usecase

import (
	"context"
	"slices"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type UpdateMemberGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewUpdateMemberGroup(groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository) *UpdateMemberGroup {
	return &UpdateMemberGroup{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}
func (u *UpdateMemberGroup) Execute(ctx context.Context, req *req.UpdateMemberGroupRequest) (*res.FailedMember, error) {
	data, err := u.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if data == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if data.ID.Hex() != req.GroupID {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	dataAction, err := u.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserActionID, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataAction == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	exists := slices.Contains(data.Settings.WhoCanApproveMember, &dataAction.Role)
	if !exists {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       *req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	datauserupdate, err := u.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, *req.UserID, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       *req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datauserupdate == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       *req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	mapper.UpdateToEntityGroupMember(req.UpdateGroupMemberReq, datauserupdate)
	err = u.groupMemberRepo.UpdateGroupMember(ctx, datauserupdate)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       *req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedMember{
		GroupID:      req.GroupID,
		UserID:       req.UserActionID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
