package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type DeleteGroupEvent struct {
	groupRepo      IRepositoryMongodb.IGroupRepository
	groupMember    IRepositoryMongodb.IGroupMembersRepository
	groupEventrepo IRepositoryMongodb.IGroupEventsRepository
}

func NewDeleteGroupEvent(groupRepo IRepositoryMongodb.IGroupRepository, groupMember IRepositoryMongodb.IGroupMembersRepository, groupEventrepo IRepositoryMongodb.IGroupEventsRepository) *DeleteGroupEvent {
	return &DeleteGroupEvent{
		groupRepo:      groupRepo,
		groupMember:    groupMember,
		groupEventrepo: groupEventrepo,
	}
}

func (uc *DeleteGroupEvent) Execute(ctx context.Context, req *req.DeleteGroupEventRequest) (*res.FailedGroupEvent, error) {
	dataGroup, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataGroup == nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	dataCreatorID, err := uc.groupMember.GetGroupMemberByUserIDAndGroupID(ctx, req.CreatorID, req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataCreatorID == nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	dataUser, err := uc.groupMember.GetGroupMemberByUserIDAndGroupID(ctx, req.UserActionID, req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if dataUser == nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if req.UserActionID != req.CreatorID {
		if dataUser.Role != sharedEnums.RoleTypeAdmin && dataUser.Role != sharedEnums.RoleTypeModerator {
			return &res.FailedGroupEvent{
				GroupID:      req.GroupID,
				EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
				ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
			}, nil
		}
	}
	err = uc.groupEventrepo.DeleteGroupEvent(ctx, req.EventID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,         // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedGroupEvent{
		GroupID:      req.GroupID,
		EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
		ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
