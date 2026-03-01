package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type AddMemberGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewAddMemberGroup(groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository) *AddMemberGroup {
	return &AddMemberGroup{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

func (uc *AddMemberGroup) Execute(ctx context.Context, req *req.AddMemberGroupRequest) (*res.FailedMember, error) {
	// Thực hiện logic thêm thành viên vào nhóm
	// Trả về lỗi nếu có

	return &res.FailedMember{
		GroupID:      req.GroupID,
		UserID:       req.UserID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
