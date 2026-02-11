package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/middleware"
	server "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/socket"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity"
	"github.com/gin-gonic/gin"
)

func NewGinServer(authMiddleware *middleware.AuthMiddleware,
	identityMod *identity.ModuleIdentity, WSHandler *server.RealtimeHandler) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/ws", WSHandler.HandleWS)
	// publicRoute := r.Group("/api/v1")
	// {

	// }
	// privateRoute := r.Group("/api/v1")
	// {

	// }
	return r
}
