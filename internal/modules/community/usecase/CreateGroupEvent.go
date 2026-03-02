package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type CreateGroupEvent struct {
	groupRepo      IRepositoryMongodb.IGroupRepository
	groupMember    IRepositoryMongodb.IGroupMembersRepository
	groupEventrepo IRepositoryMongodb.IGroupEventsRepository
}

func NewCreateGroupEvent(groupRepo IRepositoryMongodb.IGroupRepository, groupMember IRepositoryMongodb.IGroupMembersRepository, groupEventrepo IRepositoryMongodb.IGroupEventsRepository) *CreateGroupEvent {
	return &CreateGroupEvent{
		groupRepo:      groupRepo,
		groupMember:    groupMember,
		groupEventrepo: groupEventrepo,
	}
}
func (c *CreateGroupEvent) Execute(ctx context.Context, req *req.CreateGroupEventRequest) (*res.FailedGroupEvent, error) {
	datagroup, err := c.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != req.GroupID {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	datauser, err := c.groupMember.GetGroupMemberByUserIDAndGroupID(ctx, req.CreatorID, req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datauser == nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown",                    // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: errors.New("user not found"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datauser.UserID != req.CreatorID {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	entity := mapper.ToEntityGroupEvent(req.CreateGroupEventReq)
	err = c.groupEventrepo.CreateGroupEvent(ctx, entity)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedGroupEvent{
		GroupID:      req.GroupID,
		EventID:      entity.ID.Hex(), // Thay thế bằng ID thực tế của sự kiện nếu có
		ErrorMessage: nil,             // Thay thế bằng lỗi thực tế nếu có
	}, nil

}
