package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IProducer/IProducerMedia"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
)

type DeletePostUseCase struct {
	mediaAssetProducer      IProducerMedia.IProducerMedia
	pool                    IRepositoryShare.IWorkerPool
	handlerDeleteStrategies map[string]IStrategy.IPublishDeleteStrategy
}

func NewDeletePostUseCase(pool IRepositoryShare.IWorkerPool, mediaAssetProducer IProducerMedia.IProducerMedia, handlerDeleteStrategies []IStrategy.IPublishDeleteStrategy) *DeletePostUseCase {
	handlerMap := make(map[string]IStrategy.IPublishDeleteStrategy)
	for _, handler := range handlerDeleteStrategies {
		handlerMap[handler.GetDeleteType().String()] = handler
	}
	return &DeletePostUseCase{
		mediaAssetProducer:      mediaAssetProducer,
		pool:                    pool,
		handlerDeleteStrategies: handlerMap,
	}
}
func (uc *DeletePostUseCase) Execute(ctx context.Context, req *req.DeletePostRequest) (*response.Response, error) {
	callDeleteMediaAsset := func() error {
		data := mediaInContent.DeleteMediaAssetsPayloadRequestToPayload(req)
		err := uc.mediaAssetProducer.ProducerPublishPostDeleteMediaAssets(ctx, data)
		if err != nil {
			return err
		}
		return nil
	}
	resultChan := make(chan error, len(uc.handlerDeleteStrategies))
	for _, handler := range uc.handlerDeleteStrategies {
		err := uc.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			resultChan <- handler.HandlePublishDelete(safectx, req)
		})
		if err != nil {
			return response.NewResponse(response.WithData(""),
				response.WithMessage("Failed to run handler"),
				response.WithStatus(http.StatusBadRequest)), err
		}
	}
	for i := 0; i < len(uc.handlerDeleteStrategies); i++ {
		if err := <-resultChan; err != nil {
			return response.NewResponse(response.WithData(""),
				response.WithMessage("Failed to handle delete"),
				response.WithStatus(http.StatusInternalServerError)), err
		}
	}
	err := callDeleteMediaAsset()
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("Failed to delete media assets"),
			response.WithStatus(http.StatusInternalServerError)), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Post deleted successfully"), response.WithStatus(http.StatusOK)), nil
}
