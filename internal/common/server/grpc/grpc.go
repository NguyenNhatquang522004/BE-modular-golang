package grpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Engine *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	s := grpc.NewServer()
	// Best practice: Enable reflection để debug bằng Postman/gRPCui
	reflection.Register(s)
	return &GRPCServer{Engine: s}
}

func (s *GRPCServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	return s.Engine.Serve(lis)
}
