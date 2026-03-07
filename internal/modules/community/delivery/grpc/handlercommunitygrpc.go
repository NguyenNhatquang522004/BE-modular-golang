package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type HandlerCommunityGRPC struct {
	pb.UnimplementedCommunityServiceServer
	memberRepo IRepositoryMongodb.IGroupMembersRepository
}

func NewHandlerCommunityGRPC(memberRepo IRepositoryMongodb.IGroupMembersRepository) *HandlerCommunityGRPC {
	return &HandlerCommunityGRPC{
		memberRepo: memberRepo,
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
