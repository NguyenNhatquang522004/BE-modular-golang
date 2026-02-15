package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// KafkaConfig cấu hình từ env

type KafkaEventBus struct {
	writer *kafka.Writer
	config *configs.KafkaConfig
}

// NewKafkaEventBus Provider cho Wire
func NewKafkaEventBus(cfg *configs.KafkaConfig) events.EventBus {
	// Setup Producer (Writer)
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.BROKERS...),
		Balancer: &kafka.LeastBytes{},
		// Best practice: Async write để không block request chính, nhưng cần handle error
		Async: false,
		// Batch setting để tối ưu performance cho Social Network
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
	}

	return &KafkaEventBus{
		writer: w,
		config: cfg,
	}
}

func (k *KafkaEventBus) Publish(ctx context.Context, topic string, key string, eventType string, payload events.EventPayload) error {
	eventID := uuid.New().String()
	msgKey := key
	if msgKey == "" {
		msgKey = eventID
	}
	event := events.IntegrationEvent{
		ID:        eventID,
		Type:      eventType,
		Topic:     topic,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to Kafka
	err = k.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(msgKey), // KAFKA SẼ DÙNG CÁI NÀY ĐỂ CHỌN PARTITION
			Value: eventBytes,
			Topic: topic,
		},
	)
	if err != nil {
		return fmt.Errorf("kafka publish error: %w", err)
	}
	return nil
}

func (k *KafkaEventBus) Subscribe(ctx context.Context, topic string, handler events.EventHandler) error {
	// Setup Consumer (Reader) với GroupID
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.config.BROKERS,
		GroupID:  k.config.CONSUMER_GROUP, // BẮT BUỘC để scale nhiều instances
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// Run consumer trong goroutine riêng biệt (Non-blocking)
	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// Handle shutdown hoặc error connection
				return
			}

			var event events.IntegrationEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshal event: %v", err)
				continue
			}

			// Thực thi logic của module
			if err := handler(ctx, event); err != nil {
				log.Printf("Error processing event: %v", err)
				// Best practice: Có thể implement Retry hoặc Dead Letter Queue (DLQ) ở đây
			}
		}
	}()

	return nil
}

func (k *KafkaEventBus) Close() error {
	return k.writer.Close()
}
