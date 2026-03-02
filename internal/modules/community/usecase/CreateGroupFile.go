package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type CreateGroupFile struct {
	groupRepo     IRepositoryMongodb.IGroupRepository
	groupfileRepo IRepositoryMongodb.IGroupFilesRepository
}

func NewCreateGroupFile(groupRepo IRepositoryMongodb.IGroupRepository, groupfileRepo IRepositoryMongodb.IGroupFilesRepository) *CreateGroupFile {
	return &CreateGroupFile{
		groupRepo:     groupRepo,
		groupfileRepo: groupfileRepo,
	}
}

func (c *CreateGroupFile) Execute(ctx context.Context, req *req.CreateGroupFileRequest) (*res.FailGroupFile, error) {
	// Thực hiện logic tạo file cho nhóm tại đây
	// Trả về lỗi nếu có vấn đề xảy ra
	datagroup, err := c.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			ErrorMessage: err,
		}, err
	}
	if datagroup == nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			ErrorMessage: err,
		}, err
	}
	entity := mapper.ToEntityGroupFile(req.CreateGroupFileReq)
	err = c.groupfileRepo.CreateGroupFile(ctx, entity)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			ErrorMessage: err,
		}, err
	}
	return &res.FailGroupFile{
		GroupID:      req.GroupID,
		ErrorMessage: nil,
	}, nil
}
