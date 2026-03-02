package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type DeleteGroupFile struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupfileRepo   IRepositoryMongodb.IGroupFilesRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	seaweedfsRepo   IRepositoryShare.ISeaweedfs
}

func NewDeleteGroupFile(groupRepo IRepositoryMongodb.IGroupRepository, groupfileRepo IRepositoryMongodb.IGroupFilesRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, seaweedfsRepo IRepositoryShare.ISeaweedfs) *DeleteGroupFile {
	return &DeleteGroupFile{
		groupRepo:       groupRepo,
		groupfileRepo:   groupfileRepo,
		groupMemberRepo: groupMemberRepo,
		seaweedfsRepo:   seaweedfsRepo,
	}
}
func (c *DeleteGroupFile) Execute(ctx context.Context, req *req.DeleteGroupFileRequest) (*res.FailGroupFile, error) {
	// Thực hiện logic xóa file của nhóm tại đây
	// Trả về lỗi nếu có vấn đề xảy ra
	datagroup, err := c.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	if datagroup == nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	datauser, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserActionID, req.GroupID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	if datauser == nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	datafile, err := c.groupfileRepo.GetGroupFileByID(ctx, req.FileID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	if datafile == nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	if datauser.ID.Hex() != req.UserActionID {
		if datauser.Role != sharedEnums.RoleTypeAdmin && datauser.Role != sharedEnums.RoleTypeModerator {
			return &res.FailGroupFile{
				GroupID:      req.GroupID,
				UserID:       req.UserActionID,
				UserActionID: req.UserActionID,
				FileID:       req.FileID,
				ErrorMessage: errors.New("user does not have permission to delete this file"), // Thay thế bằng lỗi thực tế nếu có
			}, nil
		}
	}
	err = c.groupfileRepo.DeleteGroupFile(ctx, req.FileID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	err = c.seaweedfsRepo.Delete(ctx, datafile.StorageFileID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	return &res.FailGroupFile{
		GroupID:      req.GroupID,
		UserID:       req.UserActionID,
		UserActionID: req.UserActionID,
		FileID:       req.FileID,
		ErrorMessage: nil,
	}, nil
}
