package content

import "context"

type IProducerContent interface {
	// Các phương thức liên quan đến việc sản xuất nội dung, ví dụ:
	PublishContentPostStats(ctx context.Context) error
}
