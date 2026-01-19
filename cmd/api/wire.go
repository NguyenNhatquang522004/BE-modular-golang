//go:build wireinject
// +build wireinject

package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/google/wire"
)

func InitializeApp(config *configs.Config) (*App, error) {
	wire.Build(
		database.ConnectPostgres,
		NewGinServer,
		NewApp,
	)
	return &App{}, nil
}
