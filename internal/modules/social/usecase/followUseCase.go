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
	"github.com/google/uuid"
)

type FollowUseCase struct {
	followRepo     IRepositoryPostgres.IFollowersRepository
	followProducer IGraph.IFollowMessage
}

func NewFollowUseCase(followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage) *FollowUseCase {
	return &FollowUseCase{
		followRepo:     followRepo,
		followProducer: followProducer,
	}
}
func (f *FollowUseCase) CreateFollowUserUseCase(req *req.FollowCreateRequest) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := f.followRepo.CreateFollowUser(uuid.MustParse(req.FollowerUserID), uuid.MustParse(req.FollowedUserID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payload := &socialEvent.FollowUpdatePayload{
		Follower_UserID: req.FollowerUserID,
		Followed_UserID: req.FollowedUserID,
		IsMuted:         false,
	}
	err = f.followProducer.PublishFollowUpdateMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (f *FollowUseCase) DeleteSoftFollowUserUseCase(follower *req.FollowDeleteSoftRequest) (*response.Response, error) {
	err := f.followRepo.DeleteSoftFollowUser(uuid.MustParse(follower.FollowerUserID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := &socialEvent.FollowDeletePayload{
		Follower_UserID: follower.FollowerUserID,
	}
	err = f.followProducer.PublishFollowDeleteMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (f *FollowUseCase) DeleteHardFollowUserUseCase(follower *req.FollowDeleteHardRequest) (*response.Response, error) {
	err := f.followRepo.DeleteHardFollowUser(uuid.MustParse(follower.FollowerUserID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := &socialEvent.FollowDeletePayload{
		Follower_UserID: follower.FollowerUserID,
	}
	err = f.followProducer.PublishFollowDeleteMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (f *FollowUseCase) UpdatateMuteFollowUserUseCase(follower *req.FollowUpdateMuteRequest) (*response.Response, error) {
	err := f.followRepo.UpdatateMuteFollowUser(uuid.MustParse(follower.FollowerUserID), uuid.MustParse(follower.FollowedUserID), follower.IsMuted)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := &socialEvent.FollowUpdatePayload{
		Follower_UserID: follower.FollowerUserID,
		Followed_UserID: follower.FollowedUserID,
		IsMuted:         follower.IsMuted,
	}
	err = f.followProducer.PublishFollowUpdateMessage(ctx, payload)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (f *FollowUseCase) PaginationFollowersUseCase(req *req.FollowPaginationRequest) (*response.Response, error) {
	res, err := f.followRepo.PaginationFollowers(uuid.MustParse(req.UserID), req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(&dto.PaginationRes{
		NextCursor: res.NextCursor,
		HasNext:    res.HasNext,
		Data:       res.Data,
		Limit:      req.Limit,
	}), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (f *FollowUseCase) PaginationFollowedsUseCase(req *req.FollowPaginationRequest) (*response.Response, error) {
	res, err := f.followRepo.PaginationFolloweds(uuid.MustParse(req.UserID), req.Cursor, req.Limit)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(&dto.PaginationRes{
		NextCursor: res.NextCursor,
		HasNext:    res.HasNext,
		Data:       res.Data,
		Limit:      req.Limit,
	}), response.WithMessage("Success"), response.WithStatus("200")), nil
}
