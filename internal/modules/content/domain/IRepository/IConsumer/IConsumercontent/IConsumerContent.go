package IConsumercontent

import "context"

type IContentConsumer interface {
	// Các phương thức liên quan đến việc tiêu thụ nội dung, ví dụ:
	ConsumeContentPostStats(ctx context.Context) error
	ConsumeContentDeletePublishPost(ctx context.Context) error
}
