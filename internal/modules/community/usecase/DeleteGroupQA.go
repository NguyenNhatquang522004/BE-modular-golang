package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type DeleteGroupQA struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupqaRepo     IRepositoryMongodb.IGroupjoinQuestionsRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewDeleteGroupQA() *DeleteGroupQA {
	return &DeleteGroupQA{}
}
func (u *DeleteGroupQA) Execute(ctx context.Context, req *req.DeleteGroupQARequest) (*res.FailedGroupQA, error) {
	datagroup, err := u.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != req.GroupID {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	dataaction, err := u.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserActionID, datagroup.ID.Hex())
	if err != nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataaction == nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataaction.Role != sharedEnums.RoleTypeAdmin && dataaction.Role != sharedEnums.RoleTypeModerator {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	err = u.groupqaRepo.DeleteGroupJoinQuestion(ctx, req.QAID)
	if err != nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedGroupQA{
		GroupID:      req.GroupID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
