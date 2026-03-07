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
	data, err := h.pageRoleRepo.GetPageRoleByPageIDAndUserID(ctx, req.PageId, req.UserId)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &pb.GetRoleUserInPageResponse{
		Role: data.Role.String(),
	}, nil
}
