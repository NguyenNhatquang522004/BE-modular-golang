package IConsumerContent

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent/mediaInContent"
)

type IConsumerContent interface {
	// Định nghĩa các phương thức mà ConsumerContent sẽ thực hiện
	// Ví dụ:
	// ConsumePostCreateMediaAssets(ctx context.Context, payload *mediaInContent.CreateMediaAssetsPayload) error
	ConsumerPublishPostCreateMediaAssets(ctx context.Context, payload []*mediaInContent.CreateMediaAssetsPayload) error
	ConsumerPublishPostDeleteMediaAssets(ctx context.Context, payload []*mediaInContent.DeleteMediaAssetsPayload) error
}
