package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Server *gin.Engine
	DB     *gorm.DB
}

func NewApp(server *gin.Engine, DB *gorm.DB) *App {
	return &App{Server: server, DB: DB}
}

func NewGinServer() *gin.Engine {
	return gin.Default()
}
