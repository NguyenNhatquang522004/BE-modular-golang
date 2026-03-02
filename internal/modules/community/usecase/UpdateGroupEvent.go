package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type UpdateGroupEvent struct {
	groupRepo      IRepositoryMongodb.IGroupRepository
	groupMember    IRepositoryMongodb.IGroupMembersRepository
	groupEventrepo IRepositoryMongodb.IGroupEventsRepository
}

func NewUpdateGroupEvent(groupRepo IRepositoryMongodb.IGroupRepository, groupMember IRepositoryMongodb.IGroupMembersRepository, groupEventrepo IRepositoryMongodb.IGroupEventsRepository) *UpdateGroupEvent {
	return &UpdateGroupEvent{
		groupRepo:      groupRepo,
		groupMember:    groupMember,
		groupEventrepo: groupEventrepo,
	}
}
func (u *UpdateGroupEvent) Execute(ctx context.Context, req *req.UpdateGroupEventRequest) (*res.FailedGroupEvent, error) {
	datagroup, err := u.groupRepo.GetGroupByID(ctx, *req.GroupID)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: err,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != *req.GroupID {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      "unknown", // Thay thế bằng ID thực tế của sự kiện nếu có
			ErrorMessage: nil,       // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	dataEvent, errs := u.groupEventrepo.GetGroupEventByID(ctx, req.EventID)
	if errs != nil {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      req.EventID,
			ErrorMessage: errs, // Thay thế bằng lỗi thực tế nếu có
		}, nil

	}
	if dataEvent == nil {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      req.EventID,
			ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	mapper.UpdateToEntityGroupEvent(req.UpdateGroupEventReq, dataEvent)
	err = u.groupEventrepo.UpdateGroupEvent(ctx, dataEvent)
	if err != nil {
		return &res.FailedGroupEvent{
			GroupID:      *req.GroupID,
			EventID:      req.EventID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	return &res.FailedGroupEvent{
		GroupID:      *req.GroupID,
		EventID:      req.EventID, // Thay thế bằng ID thực tế của sự kiện nếu có
		ErrorMessage: nil,         // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
