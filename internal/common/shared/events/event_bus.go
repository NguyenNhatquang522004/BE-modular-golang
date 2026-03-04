package events

import (
	"context"
	"time"
)

// EventPayload là dữ liệu mang theo event
type EventPayload interface{}

// IntegrationEvent là cấu trúc chuẩn cho mọi event trong hệ thống
type IntegrationEvent struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Topic     string       `json:"topic"`
	Timestamp time.Time    `json:"timestamp"`
	Payload   EventPayload `json:"payload"`
}

// EventHandler là function để xử lý khi nhận được event
type EventHandler func(ctx context.Context, event IntegrationEvent) error
type BatchEventHandler func(ctx context.Context, events []IntegrationEvent) error

// EventBus interface (Dependency Inversion)
type EventBus interface {
	// THAY ĐỔI: Thêm tham số eventType vào đây
	Publish(ctx context.Context, topic string, key string, eventType string, payload EventPayload) error

	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	SubscribeBatch(ctx context.Context, topic string, batchSize int, batchTimeout time.Duration, handler BatchEventHandler) error
	Close() error
}
