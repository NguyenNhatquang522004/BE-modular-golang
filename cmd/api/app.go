package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/grpc"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server     *gin.Engine
	GRPCServer *grpc.GRPCServer
}

func NewApp(
	server *gin.Engine,
	GRPCServer *grpc.GRPCServer,
) *App {
	return &App{Server: server, GRPCServer: GRPCServer}
}
