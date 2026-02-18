package friendshipstrategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type FriendShipAccepted struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendShipAccepted(friendshipRepo IRepositoryPostgres.IFriendshipsRepository, followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage, friendshipProducer IGraph.IFriendshipMessage, blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage) *FriendShipAccepted {
	return &FriendShipAccepted{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
	}
}

// execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
func (h *FriendShipAccepted) Execute(req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	// Xử lý logic khi trạng thái là Accepted
	// Ví dụ: Cập nhật cơ sở dữ liệu, gửi thông báo, v.v.
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship accepted executed successfully"), response.WithStatus("200")), nil
}

// Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào
func (h *FriendShipAccepted) Type() enum.StatusFriendship {
	return enum.StatusFriendship_Accepted
}
