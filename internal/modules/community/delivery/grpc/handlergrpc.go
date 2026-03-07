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
	groupMember, err := h.memberRepo.GetGroupMembersByUserIDAndGroupID(ctx, req.GetUserId(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	res := &pb.GetRoleUserInGroupResponse{}
	for _, role := range groupMember {
		res.Role = append(res.Role, role.Role.String())
	}
	return res, nil
}
