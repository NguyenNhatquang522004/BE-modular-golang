package usecase

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IProducer/IProducerContent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IProducer/IProducerMedia"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublishPostUseCase struct {
	pool            IRepositoryShare.IWorkerPool
	handlerStrategy map[reflect.Type]IStrategy.IPublishPostStrategy
	producerContent IProducerContent.IProducerContent
	producerMedia   IProducerMedia.IProducerMedia
}

func NewPublishPostUseCase(
	pool IRepositoryShare.IWorkerPool, handlerStrategy []IStrategy.IPublishPostStrategy,
	producerContent IProducerContent.IProducerContent,
	producerMedia IProducerMedia.IProducerMedia) *PublishPostUseCase {
	hmap := make(map[reflect.Type]IStrategy.IPublishPostStrategy)
	for _, handler := range handlerStrategy {
		handlerType := handler.GetType()
		hmap[handlerType] = handler
	}
	return &PublishPostUseCase{
		pool:            pool,
		handlerStrategy: hmap,
		producerContent: producerContent,
		producerMedia:   producerMedia,
	}
}
func (uc *PublishPostUseCase) Execute(ctx context.Context, req *req.PublishPostRequest) (*response.Response, error) {
	reqType := reflect.TypeOf(req).Elem()
	numFields := reqType.NumField()

	// Buffer = số field tối đa có thể submit task
	resultChan := make(chan error, numFields)
	postID := primitive.NewObjectID()

	rollback := func() {
		uc.producerContent.PublishContentDeletePublishPost(ctx, &contentEvent.PostDeletePayload{PostID: postID.Hex()})
	}
	callNotification := func() {
		// todo
	}
	callMediaAssets := func() {
		capturedDataPostMedia := req.PostMedia
		if capturedDataPostMedia != nil {
			capturedDataPostMedia.PostID = postID.Hex() // Gán postID mới tạo vào payload media
			final := mediaInContent.CreateMediaAssetsPayloadMediaMetadataReqtoPayloads(capturedDataPostMedia, req.Post.UserID)
			var final2 = make([]*mediaInContent.CreateMediaAssetsPayload, 0)
			final2 = append(final2, final)
			uc.producerMedia.ProducerPublishPostCreateMediaAssets(ctx, postID.Hex(), final2)
		}
	}
	reqValue := reflect.ValueOf(req).Elem()
	submitErr := false

	for i := 0; i < numFields; i++ {
		field := reqValue.Field(i)

		// Field nil → optional, bỏ qua không rollback
		if field.Kind() != reflect.Ptr || field.IsNil() {
			continue
		}

		data := field.Interface()
		dataType := reflect.TypeOf(data)

		strategy, exists := uc.handlerStrategy[dataType]
		if !exists {
			// Non-nil field nhưng không có strategy → lỗi cấu hình
			submitErr = true
			break
		}

		// Capture biến tường minh để closure không bị override ở iteration tiếp theo
		capturedStrategy := strategy
		capturedData := data

		if err := uc.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			resultChan <- capturedStrategy.HandlePublishPost(safeCtx, capturedData, postID)
		}); err != nil {
			// Pool từ chối nhận task (full hoặc shutdown)
			submitErr = true
			break
		}
	}

	uc.pool.Wait() // Đợi toàn bộ task đang chạy hoàn thành
	close(resultChan)

	// Kiểm tra lỗi submit trước (break sớm) → rollback 1 lần duy nhất
	if submitErr {
		rollback()
		return response.NewResponse(
			response.WithData(""),
			response.WithMessage("Post Create Fail"),
			response.WithStatus(http.StatusBadRequest),
		), errors.New("failed to submit all tasks to worker pool")
	}

	// Kiểm tra kết quả từ các worker task
	for err := range resultChan {
		if err != nil {
			rollback()
			return response.NewResponse(
				response.WithData(""),
				response.WithMessage("Post Create Fail"),
				response.WithStatus(http.StatusBadRequest),
			), err
		}
	}
	callMediaAssets()  // Gọi sau khi chắc chắn tất cả task đã thành công, tránh gọi media nếu content thất bại
	callNotification() // todo
	return response.NewResponse(
		response.WithData(""),
		response.WithMessage("Post created successfully"),
		response.WithStatus(http.StatusCreated),
	), nil
}
