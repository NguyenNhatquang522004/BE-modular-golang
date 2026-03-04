package kafka

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
)

type KafkaDLQPublisher struct {
	eventBus events.EventBus
}

func NewKafkaDLQPublisher(bus events.EventBus) DLQPublisher {
	return &KafkaDLQPublisher{eventBus: bus}
}

func (p *KafkaDLQPublisher) PublishToDLQ(ctx context.Context, originalTopic string, event events.IntegrationEvent, errCause error) error {
	// Tự động tạo tên Topic DLQ (Ví dụ: "user_created" -> "user_created_dlq")
	// Bắn lại event này vào topic DLQ
	// (Tuỳ hệ thống, bạn có thể tạo 1 struct Payload mới bọc lại lỗi,
	// ở đây mình bắn lại payload cũ kèm type là DLQ_FAILED)
	dlqResult := events.NewFailedResult[any](
		event,
		errCause,
		3, // Số lần Retry tối đa (có thể lấy từ Config truyền vào)
		0, // Duration: Đoạn này ko cần thiết đo thời gian nữa, để 0
	)
	dlqResult.Data = event.Payload
	return p.eventBus.Publish(ctx, constants.TopicDLQ.String(), event.ID, "DLQ_FAILED", event.Payload)
}
