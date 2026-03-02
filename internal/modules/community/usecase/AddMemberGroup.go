package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type AddMemberGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	events          events.EventBus
}

func NewAddMemberGroup(groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, events events.EventBus) *AddMemberGroup {
	return &AddMemberGroup{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		events:          events,
	}
}

func (uc *AddMemberGroup) Execute(ctx context.Context, req *req.AddMemberGroupRequest) (*res.FailedMember, error) {
	// Thực hiện logic thêm thành viên vào nhóm
	// Trả về lỗi nếu có
	datagroup, err := uc.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group not found"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	if datagroup.ID.Hex() != req.GroupID {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: errors.New("group ID mismatch"), // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	entitymember := mapper.ToEntityGroupMember(req.CreateGroupMemberReq)
	entitymember.Role = sharedEnums.RoleTypeMember
	entitymember.Status = sharedEnums.ProcessingPending
	err = uc.groupMemberRepo.CreateGroupMember(ctx, entitymember)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			ErrorMessage: err, // Thay thế bằng lỗi thực tế nếu có
		}, nil
	}
	payload := &communityEvent.GroupStatsPayload{
		GroupID:            req.GroupID,
		MemberCount:        0,
		PostCount:          0,
		PendingMemberCount: 1,
		PendingPostCount:   0,
		ReportedPostCount:  0,
		EventType:          constants.Created,
	}
	err = uc.events.Publish(ctx, constants.TopicGroupStats.String(), datagroup.ID.Hex(), constants.Created.String(), payload)
	return &res.FailedMember{
		GroupID:      req.GroupID,
		UserID:       req.UserID,
		ErrorMessage: nil, // Thay thế bằng lỗi thực tế nếu có
	}, nil
}
