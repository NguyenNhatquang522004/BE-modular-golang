package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
)

type ICreateAlbumUseCase interface {
	Execute(ctx context.Context, req *req.CreateAlbumRequest) (*res.CreateAlbumResponse, error)
}
type IUpdateAlbumUseCase interface {
	Execute(ctx context.Context, req *req.UpdateAlbumRequest) (*res.UpdateAlbumResponse, error)
}
type IDeleteAlbumUseCase interface {
	Execute(ctx context.Context, req *req.DeleteAlbumRequest) (*res.DeleteAlbumResponse, error)
}
type IReactAlbumUseCase interface {
	Execute(ctx context.Context, req *req.ReactAlbumRequest) (*res.ReactAlbumResponse, error)
}
type ICreateStoryUseCase interface {
	Execute(ctx context.Context, req *req.CreateStoryRequest) (*res.FailedStoryResponse, error)
}
type IUpdateStoryUseCase interface {
	Execute(ctx context.Context, req *req.UpdateStoryRequest) (*res.FailedStoryResponse, error)
}
type IDeleteStoryUseCase interface {
	Execute(ctx context.Context, req *req.DeleteStoryRequest) (*res.FailedStoryResponse, error)
}
type IViewCountStoryUseCase interface {
	Execute(ctx context.Context, req *req.ViewCountStoryRequest) (*res.FailedStoryResponse, error)
}
type IReactStoryUseCase interface {
	Execute(ctx context.Context, req *req.ReactStoryRequest) (*res.FailedStoryResponse, error)
}
type IRelyStoryUseCase interface {
	Execute(ctx context.Context, req *req.RelyStoryRequest) (*[]res.FailedStoryResponse, error)
}
type ICreateReelUseCase interface {
	Execute(ctx context.Context, req *req.CreateReelRequest) (*res.FailedReelResponse, error)
}
type IUpdateReelUseCase interface {
	Execute(ctx context.Context, req *req.UpdateReelRequest) (*res.FailedReelResponse, error)
}
type IDeleteReelUseCase interface {
	Execute(ctx context.Context, req *req.DeleteReelRequest) (*res.FailedReelResponse, error)
}
type IReactReelUseCase interface {
	Execute(ctx context.Context, req *req.ReactReelRequest) (*res.FailedReelResponse, error)
}
type IReactCounterReelUseCase interface {
	Execute(ctx context.Context, req *req.ReactCounterReelRequest) (*res.FailedReelResponse, error)
}
type IShareReelUseCase interface {
	Execute(ctx context.Context, req *req.ShareReelRequest) (*res.FailedReelResponse, error)
}
type ICreateLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.CreateLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type IUpdateLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.UpdateLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type IDeleteLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.DeleteLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type IReactLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.ReactLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type ICounterLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.CounterLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type IStartStopVideoLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.StartStopVideoLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type ICommentLiveStreamUseCase interface {
	Execute(ctx context.Context, req *req.CommentLiveStreamRequest) (*res.FailedLiveStreamResponse, error)
}
type ICreateArtistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IUpdateArtistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IDeleteArtistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IReactArtistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type ICreateItemMusicLibraryUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IDeleteItemMusicLibraryUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IReactItemMusicLibraryUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type ICreatePlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IUpdatePlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IDeletePlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IReactPlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IAddItemToPlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type IRemoveItemFromPlaylistUseCase interface {
	Execute(ctx context.Context) (*response.Response, error)
}
type usecase struct {
}

func NewUsecase() *usecase {
	return &usecase{}
}
