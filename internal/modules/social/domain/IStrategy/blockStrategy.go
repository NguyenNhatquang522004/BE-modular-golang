package IStrategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type IBlockStrategy interface {
	Execute(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error)
	Type() enum.Type_Block // Để nhận diện strategy này dùng cho Enum nào
}
