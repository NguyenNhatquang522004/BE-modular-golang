package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type RemoveMemberGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	groupfileRepo   IRepositoryMongodb.IGroupFilesRepository
}

func NewRemoveMemberGroup() *RemoveMemberGroup {
	return &RemoveMemberGroup{}
}
func (uc *RemoveMemberGroup) Execute(ctx context.Context, req *req.RemoveMemberGroupRequest) (*res.FailedMember, error) {
	// Thực hiện logic xóa thành viên khỏi nhóm
	// Trả về lỗi nếu có
	datagroup, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group not found"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != req.GroupID {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group ID mismatch"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	data, err := uc.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserID, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err,
		}, err
	}
	if data == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("Member not found in the group"),
		}, nil
	}
	dataAction, err := uc.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.GroupID, req.UserActionID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err,
		}, err
	}
	if dataAction == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("Unauthorized: User performing the action is not a member of the group"),
		}, nil
	}
	// Kiểm tra quyền hạn của người thực hiện hành động (admin hoặc chính user đó)
	if dataAction.Role != sharedEnums.RoleTypeAdmin && req.UserActionID != req.UserID {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("Unauthorized: Only admins or the user themselves can remove a member"),
		}, nil
	}
	datauser, err := uc.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.GroupID, req.UserID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err,
		}, err
	}
	if datauser == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("Member not found in the group"),
		}, nil
	}
	if datauser.Role == sharedEnums.RoleTypeAdmin {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("Cannot remove an admin member from the group"),
		}, nil
	}
	err = uc.groupMemberRepo.DeleteGroupMemberByUserIDAndGroupID(ctx, req.GroupID, req.UserID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err,
		}, err
	}
	err = uc.groupfileRepo.DeleteGroupFilesByGroupIDAndUploaderID(ctx, data.ID.Hex(), req.UserID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err,
		}, err
	}
	return &res.FailedMember{
		GroupID:      req.GroupID,
		UserID:       req.UserID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
