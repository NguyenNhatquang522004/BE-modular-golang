//go:build wireinject
// +build wireinject

package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/grpc"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/client"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity"
	"github.com/google/wire"
)

func InitializeApp(config *configs.Config) (*App, func(), error) {
	wire.Build(
		client.ProviderSet,
		identity.ModuleIndentitySet,
		grpc.NewGRPCServer,
		NewGinServer,
		NewApp,
	)
	return &App{}, nil, nil
}
