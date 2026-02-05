package gateway

import (
	v1 "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/gateway/delivery/http/v1"
	"github.com/google/wire"
)

var routeGateWay = wire.NewSet(
	v1.NewRouterV1,
)
var middlewareGateWay = wire.NewSet(
	routeGateWay,
)
var ModuleGateway = wire.NewSet(
	routeGateWay,
	middlewareGateWay,
)	
