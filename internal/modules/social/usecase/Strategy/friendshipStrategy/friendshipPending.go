package friendshipstrategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type FriendshipPending struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendshipPending(friendshipRepo IRepositoryPostgres.IFriendshipsRepository, followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage, friendshipProducer IGraph.IFriendshipMessage, blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage) *FriendshipPending {
	return &FriendshipPending{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
	}
}

// Execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
func (h *FriendshipPending) Execute(req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	// Xử lý logic khi trạng thái là Pending
	// Ví dụ: Cập nhật cơ sở dữ liệu, gửi thông báo, v.v.
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship pending executed successfully"), response.WithStatus("200")), nil
}

// Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào

func (h *FriendshipPending) Type() enum.StatusFriendship {
	return enum.StatusFriendship_Pending
}
