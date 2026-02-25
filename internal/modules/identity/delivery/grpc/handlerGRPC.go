package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IdentityHandlerGRPC struct {
	pb.UnimplementedIdentityServiceServer
}

func NewIdentityHandlerGRPC() *IdentityHandlerGRPC {
	return &IdentityHandlerGRPC{}
}

func (h *IdentityHandlerGRPC) GetUserOTP(ctx context.Context, req *emptypb.Empty) (*pb.UserIDResponse, error) {
	// 1. Bóc Metadata từ Context (Incoming)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Thiếu metadata")
	}

	// 2. Lấy giá trị động ra bằng cái Key mà A đã nhét vào
	userIDs := md["x-user-id"]
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Không tìm thấy x-user-id trong context")
	}

	// Đây chính là dữ liệu "động" mà A truyền sang!
	dynamicID := userIDs[0]

	// 3. Bây giờ B có ID động rồi, B có thể gọi Database để lấy data thật
	// user, err := h.usecase.FindUserInDB(dynamicID)

	// Trả về cho A
	return &pb.UserIDResponse{
		UserId: dynamicID, // Trả lại đúng cái ID động đó
	}, nil
}
