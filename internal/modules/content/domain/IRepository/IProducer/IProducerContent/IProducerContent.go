package IProducerContent

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
)

type IProducerContent interface {
	// Các phương thức liên quan đến việc sản xuất nội dung, ví dụ:
	PublishContentPostStats(ctx context.Context, payload ...*contentEvent.PostStatsPayload) error
	PublishContentDeletePublishPost(ctx context.Context, payload ...*contentEvent.PostDeletePayload) error
}
