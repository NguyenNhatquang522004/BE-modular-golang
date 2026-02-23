package usecase

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IProducer/IProducerContent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublishPostUseCase struct {
	pool            IRepositoryShare.IWorkerPool
	handlerStrategy map[reflect.Type]IStrategy.IPublishPostStrategy
	producerContent IProducerContent.IProducerContent
}

func NewPublishPostUseCase(
	pool IRepositoryShare.IWorkerPool, handlerStrategy []IStrategy.IPublishPostStrategy, producerContent IProducerContent.IProducerContent) *PublishPostUseCase {
	hmap := make(map[reflect.Type]IStrategy.IPublishPostStrategy)
	for _, handler := range handlerStrategy {
		handlerType := handler.GetType()
		hmap[handlerType] = handler
	}
	return &PublishPostUseCase{
		pool:            pool,
		handlerStrategy: hmap,
		producerContent: producerContent,
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

	return response.NewResponse(
		response.WithData(""),
		response.WithMessage("Post created successfully"),
		response.WithStatus(http.StatusCreated),
	), nil
}
func (uc *PublishPostUseCase) ExecuteBulk(ctx context.Context, req []*req.PublishPostRequest) (*response.Response, error) {

	return nil, nil
}
