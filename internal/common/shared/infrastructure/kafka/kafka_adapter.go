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
func (k *KafkaEventBus) SubscribeBatch(ctx context.Context, topic string, batchSize int, batchTimeout time.Duration, handler events.BatchEventHandler) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.config.BROKERS,
		GroupID:  k.config.CONSUMER_GROUP,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		defer r.Close()

		// Channel để nhận message liên tục từ Kafka
		msgChan := make(chan kafka.Message)

		// Goroutine phụ: Liên tục Fetch message (Non-blocking commit) đẩy vào channel
		go func() {
			for {
				m, err := r.FetchMessage(ctx)
				if err != nil {
					// Xử lý khi context bị hủy hoặc reader đóng
					log.Printf("Fetch message error or context cancelled: %v", err)
					close(msgChan)
					return
				}
				msgChan <- m
			}
		}()

		var batchMsgs []kafka.Message
		var batchEvents []events.IntegrationEvent

		// Ticker để xả (flush) batch nếu đợi quá lâu mà chưa đủ batchSize
		ticker := time.NewTicker(batchTimeout)
		defer ticker.Stop()

		// Hàm helper để xử lý và commit
		processAndCommit := func() {
			if len(batchEvents) == 0 {
				return
			}

			// 1. Gọi handler xử lý logic (DB insert, bulk update...)
			err := handler(ctx, batchEvents)
			if err != nil {
				log.Printf("Error processing batch: %v", err)
				// Lưu ý: Cần có chiến lược Retry hoặc DLQ ở đây.
				// Nếu không commit, Kafka sẽ gửi lại các message này khi restart.
				return
			}

			// 2. Nếu xử lý thành công, Commit toàn bộ message
			if err := r.CommitMessages(ctx, batchMsgs...); err != nil {
				log.Printf("Error committing messages: %v", err)
			}

			// 3. Reset lại batch
			batchMsgs = batchMsgs[:0]
			batchEvents = batchEvents[:0]
		}

		// Vòng lặp chính gom batch
		for {
			select {
			case m, ok := <-msgChan:
				if !ok {
					return // Channel bị đóng
				}

				var event events.IntegrationEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					log.Printf("Error unmarshal event: %v", err)
					// Bỏ qua message lỗi format nhưng CẦN commit nó để không bị kẹt
					r.CommitMessages(ctx, m)
					continue
				}

				batchMsgs = append(batchMsgs, m)
				batchEvents = append(batchEvents, event)

				// Nếu đủ size thì xử lý ngay
				if len(batchEvents) >= batchSize {
					processAndCommit()
					ticker.Reset(batchTimeout) // Reset lại đồng hồ đếm ngược
				}

			case <-ticker.C:
				// Hết thời gian chờ, có bao nhiêu xử lý bấy nhiêu
				processAndCommit()

			case <-ctx.Done():
				// Ứng dụng shutdown, cố gắng xử lý nốt batch còn tồn đọng
				processAndCommit()
				return
			}
		}
	}()

	return nil
}
func (k *KafkaEventBus) Close() error {
	return k.writer.Close()
}
