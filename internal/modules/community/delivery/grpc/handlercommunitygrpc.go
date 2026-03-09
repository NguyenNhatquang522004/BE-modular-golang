package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type HandlerCommunityGRPC struct {
	pb.UnimplementedCommunityServiceServer
	memberRepo IRepositoryMongodb.IGroupMembersRepository
	groupRepo  IRepositoryMongodb.IGroupRepository
}

func NewHandlerCommunityGRPC(memberRepo IRepositoryMongodb.IGroupMembersRepository, groupRepo IRepositoryMongodb.IGroupRepository) *HandlerCommunityGRPC {
	return &HandlerCommunityGRPC{
		memberRepo: memberRepo,
		groupRepo:  groupRepo,
	}
}

func (h *HandlerCommunityGRPC) GetRoleUserInGroup(ctx context.Context, req *pb.GetRoleUserInGroupRequest) (*pb.GetRoleUserInGroupResponse, error) {
	groupMember, err := h.memberRepo.GetGroupMemberByUserIDAndGroupID(ctx, req.GetUserId(), req.GetGroupId())
	if err != nil {
		return nil, err
	}

	return &pb.GetRoleUserInGroupResponse{
		Role: groupMember.Role.String(),
	}, nil
}
func (h *HandlerCommunityGRPC) GetGroupInfo(ctx context.Context, req *pb.GetGroupInfoRequest) (*pb.GetGroupInfoResponse, error) {
	data, err := h.groupRepo.GetGroupByID(ctx, req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, err
	}
	return &pb.GetGroupInfoResponse{
		GroupId:   data.ID.Hex(),
		Protected: data.Privacy.String(),
	}, nil
}
