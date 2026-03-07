package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type HandlerBusinessGRPC struct {
	pb.UnimplementedBusinessServiceServer
	pageRoleRepo IRepositoryMongodb.IPageRolesRepository
}

func NewHandlerBusinessGRPC(pageRoleRepo IRepositoryMongodb.IPageRolesRepository) *HandlerBusinessGRPC {
	return &HandlerBusinessGRPC{
		pageRoleRepo: pageRoleRepo,
	}
}
func (h *HandlerBusinessGRPC) GetRoleUserInPage(ctx context.Context, req *pb.GetRoleUserInPageRequest) (*pb.GetRoleUserInPageResponse, error) {
	data, err := h.pageRoleRepo.GetPageRolesByPageIDAndUserID(ctx, req.PageId, req.UserId)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return &pb.GetRoleUserInPageResponse{
			Role: []string{},
		}, nil
	}
	var roles []string
	for _, role := range data {
		roles = append(roles, role.Role.String())
	}
	return &pb.GetRoleUserInPageResponse{
		Role: roles,
	}, nil
}
