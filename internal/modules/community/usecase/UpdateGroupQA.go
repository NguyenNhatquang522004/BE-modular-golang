package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type UpdateGroupQA struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupqaRepo     IRepositoryMongodb.IGroupjoinQuestionsRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewUpdateGroupQA(groupRepo IRepositoryMongodb.IGroupRepository, groupqaRepo IRepositoryMongodb.IGroupjoinQuestionsRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository) *UpdateGroupQA {
	return &UpdateGroupQA{
		groupRepo:       groupRepo,
		groupqaRepo:     groupqaRepo,
		groupMemberRepo: groupMemberRepo,
	}
}
func (u *UpdateGroupQA) Execute(ctx context.Context, req *req.UpdateGroupQARequest) (*res.FailedGroupQA, error) {
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
	dataQA, err := u.groupqaRepo.GetGroupJoinQuestionByID(ctx, req.QAID)
	if err != nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataQA == nil {
		return &res.FailedGroupQA{
			GroupID:      req.GroupID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	mapper.UpdateToEntityGroupJoinQuestion(req.UpdateGroupJoinQuestionReq, dataQA)
	err = u.groupqaRepo.UpdateGroupJoinQuestion(ctx, dataQA)
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
