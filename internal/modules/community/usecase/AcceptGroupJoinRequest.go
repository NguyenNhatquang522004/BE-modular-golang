package usecase

import (
	"context"
	"slices"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type AcceptGroupJoinRequest struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	events          events.EventBus
}

func NewAcceptGroupJoinRequest(groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, events events.EventBus) *AcceptGroupJoinRequest {
	return &AcceptGroupJoinRequest{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		events:          events,
	}
}

func (a *AcceptGroupJoinRequest) Execute(ctx context.Context, req *req.AcceptGroupJoinRequest) (*res.FailedMember, error) {
	// Implement the logic for accepting a group join request here
	datagroup, err := a.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	if datagroup == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	dataUserAction, err := a.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserActionID, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	if dataUserAction == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	exists := slices.Contains(datagroup.Settings.WhoCanApproveMember, &dataUserAction.Role)
	if !exists {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	// Tiếp tục với logic chấp nhận hoặc từ chối yêu cầu gia nhập nhóm
	datauser, err := a.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.UserID, req.GroupID)
	if err != nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	if datauser == nil {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}
	if datauser.Status != sharedEnums.ProcessingPending {
		return &res.FailedMember{
			GroupID:      req.GroupID,
			UserID:       req.UserID,
			UserActionID: req.UserActionID,
			ErrorMessage: err,
		}, err
	}

	if req.Accept { // Logic để chấp nhận yêu cầu gia nhập nhóm
		datauser.Status = sharedEnums.ProcessingActive
		err = a.groupMemberRepo.UpdateGroupMember(ctx, datauser)
		if err != nil {
			return &res.FailedMember{
				GroupID:      req.GroupID,
				UserID:       req.UserID,
				UserActionID: req.UserActionID,
				ErrorMessage: err,
			}, err
		}
		payload := &communityEvent.GroupStatsPayload{
			GroupID:            req.GroupID,
			MemberCount:        1,
			PostCount:          0,
			PendingMemberCount: -1,
			PendingPostCount:   0,
			ReportedPostCount:  0,
			EventType:          constants.Created,
		}
		err = a.events.Publish(ctx, constants.TopicGroupStats.String(), datagroup.ID.Hex(), constants.Created.String(), payload)

	} else {
		// Logic để từ chối yêu cầu gia nhập nhóm
		err = a.groupMemberRepo.DeleteGroupMember(ctx, datauser.ID.Hex())
		if err != nil {
			return &res.FailedMember{
				GroupID:      req.GroupID,
				UserID:       req.UserID,
				UserActionID: req.UserActionID,
				ErrorMessage: err,
			}, err
		}
		payload := &communityEvent.GroupStatsPayload{
			GroupID:            req.GroupID,
			MemberCount:        0,
			PostCount:          0,
			PendingMemberCount: -1,
			PendingPostCount:   0,
			ReportedPostCount:  0,
			EventType:          constants.Created,
		}
		err = a.events.Publish(ctx, constants.TopicGroupStats.String(), datagroup.ID.Hex(), constants.Created.String(), payload)
	}
	return nil, nil
}
