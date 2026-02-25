package IConsumerContent

import (
	"context"
)

type IConsumerContent interface {
	// Định nghĩa các phương thức mà ConsumerContent sẽ thực hiện
	// Ví dụ:
	// ConsumePostCreateMediaAssets(ctx context.Context, payload *mediaInContent.CreateMediaAssetsPayload) error
	ConsumerPublishPostCreateDeleteMediaAssets(ctx context.Context) error

}
