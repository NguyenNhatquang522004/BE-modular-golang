package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/grpc"
	server "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/socket"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/socket"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server       *gin.Engine
	GRPCServer   *grpc.GRPCServer
	SocketServer *server.RealtimeHandler
	Hub          socket.Manager
}

func NewApp(
	server *gin.Engine,
	GRPCServer *grpc.GRPCServer,
	SocketServer *server.RealtimeHandler,
	Hub socket.Manager,
) *App {
	return &App{Server: server, GRPCServer: GRPCServer, SocketServer: SocketServer, Hub: Hub}
}
