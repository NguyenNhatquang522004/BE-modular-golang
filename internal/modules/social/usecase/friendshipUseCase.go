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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type FriendshipUseCase struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendshipUseCase(friendshipRepo IRepositoryPostgres.IFriendshipsRepository,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IGraph.IFollowMessage,
	friendshipProducer IGraph.IFriendshipMessage,
	blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IGraph.IBlockMessage) *FriendshipUseCase {
	return &FriendshipUseCase{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
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
func (uc *FriendshipUseCase) UpdateFriendshipStatusUseCase(req *req.UpdateFriendshipStatusRequest) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch enum.StatusFriendship(req.Status) {
	case enum.StatusFriendship_Accepted:
		frienshipdata, err := uc.friendshipRepo.GetFriendshipTableByTableId(uuid.MustParse(req.FriendshipID))
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
		err = uc.followRepo.CreateFollowUser(frienshipdata.Recipient_ID, frienshipdata.Requester_ID)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
		payloadfollow := &socialEvent.FollowUpdatePayload{
			Follower_UserID: frienshipdata.Recipient_ID.String(),
			Followed_UserID: frienshipdata.Requester_ID.String(),
			IsMuted:         false,
		}
		payloadfollow1 := &socialEvent.FollowUpdatePayload{
			Follower_UserID: frienshipdata.Requester_ID.String(),
			Followed_UserID: frienshipdata.Recipient_ID.String(),
			IsMuted:         false,
		}
		err = uc.followProducer.PublishFollowUpdateMessage(ctx, payloadfollow, payloadfollow1)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
		// nếu chấp nhận lời mời kết bạn thì sẽ tạo thêm một bản ghi follow ngược lại nếu đã có rồi thì không cần tạo nữa
	case enum.StatusFriendship_Declined:
		err := uc.friendshipRepo.UpdateFriendshipStatus(uuid.MustParse(req.FriendshipID), enum.StatusFriendship(req.Status))
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
	}
	return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), nil
}

func (uc *FriendshipUseCase) UpdateBlockFriendshipUseCase(req *req.UpdateBlockFriendshipRequest) (*response.Response, error) {
	// Implement the logic for updating block friendship here
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	frienshipdata, err := uc.friendshipRepo.GetFriendshipTableByTableId(uuid.MustParse(req.FriendshipID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	err = uc.blockRepo.CreateBlockUser(frienshipdata.Requester_ID, frienshipdata.Recipient_ID, req.Status)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}

	payload := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: frienshipdata.Requester_ID.String(),
		Blocked_UserID: frienshipdata.Recipient_ID.String(),
		Status:         enum.Type_Block(req.Status),
	}
	err = uc.blockProducer.PublishBlockUpdateMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	if frienshipdata.Status == enum.StatusFriendship_Accepted {
		err := uc.friendshipRepo.UpdateFriendshipStatus(uuid.MustParse(req.FriendshipID), enum.StatusFriendship(req.Status))
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
		payload1 := &socialEvent.FollowDeletePayload{
			Follower_UserID: frienshipdata.Requester_ID.String(),
			Followed_UserID: frienshipdata.Recipient_ID.String(),
		}
		payload2 := &socialEvent.FollowDeletePayload{
			Follower_UserID: frienshipdata.Recipient_ID.String(),
			Followed_UserID: frienshipdata.Requester_ID.String(),
		}
		err = uc.followProducer.PublishFollowDeleteMessage(ctx, payload1, payload2)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}

		err = uc.followRepo.DeleteBatchSoftFollowUser(frienshipdata.Requester_ID, frienshipdata.Recipient_ID)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
		}
	}

	err = uc.friendshipRepo.UpdateFriendshipStatus(uuid.MustParse(req.FriendshipID), enum.StatusFriendship(req.Status))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	followtable1, err := uc.followRepo.GetFollowerTableBybidirectional(frienshipdata.Requester_ID, frienshipdata.Recipient_ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	payload1 := &socialEvent.FollowDeletePayload{
		Follower_UserID: followtable1.Follower_UserID.String(),
		Followed_UserID: followtable1.Followed_UserID.String(),
	}
	payload2 := &socialEvent.FollowDeletePayload{
		Follower_UserID: followtable1.Followed_UserID.String(),
		Followed_UserID: followtable1.Follower_UserID.String(),
	}
	defer cancel()
	err = uc.followProducer.PublishFollowDeleteMessage(ctx, payload1, payload2)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	err = uc.followRepo.DeleteSoftFollowUser(followtable1.ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), err
	}
	// nếu block thì sẽ xóa mềm tất cả các bản ghi follow liên quan đến 2 người này và tạo một bản ghi block trong bảng friendship còn nếu đã block rồi thì không cần tạo nữa

	return response.NewResponse(response.WithData(""), response.WithMessage(""), response.WithStatus("")), nil
}
func (uc *FriendshipUseCase) PanigationAcceptedFriendshipUseCase(req *req.PaginationFriendshipRequest) (*response.Response, error) {
	data, err := uc.friendshipRepo.PanigationAcceptedFriendship(uuid.MustParse(req.UserID), req.Metadata.Cursor, req.Metadata.Limit)
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
func (uc *FriendshipUseCase) PanigationPendingFriendshipUseCase(req *req.PaginationFriendshipRequest) (*response.Response, error) {
	data, err := uc.friendshipRepo.PanigationPendingFriendship(uuid.MustParse(req.UserID), req.Metadata.Cursor, req.Metadata.Limit)
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
