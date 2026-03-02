package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type CreateGroupQA struct {
	groupRepo   IRepositoryMongodb.IGroupRepository
	groupqaRepo IRepositoryMongodb.IGroupjoinQuestionsRepository
}

func NewCreateGroupQA(groupRepo IRepositoryMongodb.IGroupRepository, groupqaRepo IRepositoryMongodb.IGroupjoinQuestionsRepository) *CreateGroupQA {
	return &CreateGroupQA{
		groupRepo:   groupRepo,
		groupqaRepo: groupqaRepo,
	}
}
func (u *CreateGroupQA) Execute(ctx context.Context, req *req.CreateGroupQARequest) (*res.FailedGroupQA, error) {
	// Thực hiện logic tạo Group QA tại đây
	// Nếu có lỗi, trả về lỗi trong res.FailedGroupQA
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
	entityqa := mapper.ToEntityGroupJoinQuestion(req.CreateGroupJoinQuestionReq)
	err = u.groupqaRepo.CreateGroupJoinQuestion(ctx, entityqa)
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
