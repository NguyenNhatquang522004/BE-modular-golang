package IStrategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type IFriendshipStrategy interface {
	Execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
	Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào
}
