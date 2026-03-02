package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type DownloadGroupFile struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupfileRepo   IRepositoryMongodb.IGroupFilesRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	seaweedfsRepo   IRepositoryShare.ISeaweedfs
	events          events.EventBus
}

func NewDownloadGroupFile(groupRepo IRepositoryMongodb.IGroupRepository, groupfileRepo IRepositoryMongodb.IGroupFilesRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, seaweedfsRepo IRepositoryShare.ISeaweedfs, events events.EventBus) *DownloadGroupFile {
	return &DownloadGroupFile{
		groupRepo:       groupRepo,
		groupfileRepo:   groupfileRepo,
		groupMemberRepo: groupMemberRepo,
		seaweedfsRepo:   seaweedfsRepo,
		events:          events,
	}
}
func (c *DownloadGroupFile) Execute(ctx context.Context, req *req.DownloadGroupFileRequest) (*res.FailGroupFile, error) {
	// Thực hiện logic xóa file của nhóm tại đây
	// Trả về lỗi nếu
	// có vấn đề xảy ra
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
	datafile.DownloadCount += 1
	err = c.groupfileRepo.UpdateGroupFileDownloadCount(ctx, datafile.ID.Hex(), datafile.DownloadCount+1)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	pr, err := c.seaweedfsRepo.Download(ctx, datafile.StorageFileID)
	if err != nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	if pr == nil {
		return &res.FailGroupFile{
			GroupID:      req.GroupID,
			UserID:       req.UserActionID,
			UserActionID: req.UserActionID,
			FileID:       req.FileID,
			ErrorMessage: err,
		}, err
	}
	err = c.events.Publish(ctx, constants.TopicDownloadGroupFile.String(), req.GroupID, constants.Created.String(), &communityEvent.DownloadGroupFilePayload{
		GroupID:   req.GroupID,
		FileID:    req.FileID,
		Count:     1,
		UserID:    req.UserActionID,
		EventType: constants.Created,
	})
	return &res.FailGroupFile{
		GroupID:      req.GroupID,
		UserID:       req.UserActionID,
		UserActionID: req.UserActionID,
		FileID:       req.FileID,
		Data:         pr,
		ErrorMessage: nil,
	}, nil
}
