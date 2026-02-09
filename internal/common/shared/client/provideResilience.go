package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/resilience"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	resilience.NewBreakerProvider,
)
