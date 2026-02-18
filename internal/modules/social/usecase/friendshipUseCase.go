package usecase

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IStrategy"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type FriendshipUseCase struct {
	friendshipRepo            IRepositoryPostgres.IFriendshipsRepository
	followRepo                IRepositoryPostgres.IFollowersRepository
	followProducer            IGraph.IFollowMessage
	friendshipProducer        IGraph.IFriendshipMessage
	blockRepo                 IRepositoryPostgres.IBlockRepository
	blockProducer             IGraph.IBlockMessage
	handlerFriendshipStrategy map[enum.StatusFriendship]IStrategy.IFriendshipStrategy
}

func NewFriendshipUseCase(friendshipRepo IRepositoryPostgres.IFriendshipsRepository,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IGraph.IFollowMessage,
	friendshipProducer IGraph.IFriendshipMessage,
	blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IGraph.IBlockMessage, handlerFriendshipStrategy []IStrategy.IFriendshipStrategy) *FriendshipUseCase {
	Hmap := make(map[enum.StatusFriendship]IStrategy.IFriendshipStrategy)
	for _, handler := range handlerFriendshipStrategy {
		Hmap[handler.Type()] = handler
	}
	return &FriendshipUseCase{
		friendshipRepo:            friendshipRepo,
		followRepo:                followRepo,
		followProducer:            followProducer,
		friendshipProducer:        friendshipProducer,
		blockRepo:                 blockRepo,
		blockProducer:             blockProducer,
		handlerFriendshipStrategy: Hmap,
	}
}

func (uc *FriendshipUseCase) CreateFriendshipUseCase(req *req.CreateFriendshipRequest) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checkBlock, err := uc.blockRepo.IsBlocked(uuid.MustParse(req.RequesterID), uuid.MustParse(req.RecipientID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	if checkBlock {
		return response.NewResponse(response.WithData(""), response.WithMessage("You are blocked by this user or you have blocked this user"), response.WithStatus("403")), nil
	}
	err = uc.friendshipRepo.CreateFriendship(uuid.MustParse(req.RequesterID), uuid.MustParse(req.RecipientID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	payload := &socialEvent.FriendshipsUpdatePayload{
		Requester_ID: req.RequesterID,
		Recipient_ID: req.RecipientID,
		Status:       enum.StatusFriendship_Pending,
	}
	err = uc.friendshipProducer.PublishFriendshipUpdateMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	// sau này nếu trong user sẽ có phần user này đang public hay private thì sẽ check ở đây nếu private thì sẽ không tạo follow còn public thì sẽ tạo follow luôn
	// phải gọi phía module identity để check xem user có đang private hay không nếu private thì sẽ không tạo follow còn public thì sẽ tạo follow luôn bằng GRPC
	err = uc.followRepo.CreateFollowUser(uuid.MustParse(req.RequesterID), uuid.MustParse(req.RecipientID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	payloadfollow := &socialEvent.FollowUpdatePayload{
		Follower_UserID: req.RecipientID,
		Followed_UserID: req.RequesterID,
		IsMuted:         false,
	}
	err = uc.followProducer.PublishFollowUpdateMessage(ctx, payloadfollow)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship created successfully"), response.WithStatus("200")), nil
}
func (uc *FriendshipUseCase) HandleFriendShipUseCase(req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	handler, exists := uc.handlerFriendshipStrategy[req.Status]
	if !exists {
		return response.NewResponse(response.WithData(""), response.WithMessage("No handler found for the given friendship status"), response.WithStatus("400")), nil
	}
	data, err := handler.Execute(req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), err
	}
	return data, nil
}
func (uc *FriendshipUseCase) PanigationAcceptedFriendshipUseCase(req *req.PaginationFriendshipRequest) (*response.Response, error) {
	data, err := uc.friendshipRepo.PanigationStatusFriendship(uuid.MustParse(req.UserID), req.Metadata.Cursor, req.Metadata.Limit, enum.StatusFriendship_Accepted)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(&dto.PaginationRes{
		NextCursor: data.NextCursor,
		HasNext:    data.HasNext,
		Data:       data.Data,
		Limit:      data.Limit,
	}), response.WithMessage(""), response.WithStatus("")), nil
}
