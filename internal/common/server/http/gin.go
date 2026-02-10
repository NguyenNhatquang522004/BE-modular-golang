package http

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/middleware"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity"
	"github.com/gin-gonic/gin"
)

func NewGinServer(authMiddleware *middleware.AuthMiddleware,
	identityMod *identity.ModuleIdentity) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	// publicRoute := r.Group("/api/v1")
	// {

	// }
	// privateRoute := r.Group("/api/v1")
	// {

	// }
	return r
}
