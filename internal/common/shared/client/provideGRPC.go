package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	pb "github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProvideGRPCConnection(cfg *configs.Config) (*grpc.ClientConn, error) {
	// Trong môi trường Monolith, ta gọi localhost
	return grpc.Dial(cfg.GRPCServer.GRPC_SERVER_HOST+":"+cfg.GRPCServer.GRPC_SERVER_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Tắt SSL nội bộ
	)
}
func ProvideIdentityClient(conn *grpc.ClientConn) pb.IdentityServiceClient {
	return pb.NewIdentityServiceClient(conn)
}
