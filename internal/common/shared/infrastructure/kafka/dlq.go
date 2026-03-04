package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
)

// =====================================================================
// 1. ĐỊNH NGHĨA CUSTOM ERROR ĐỂ PHÂN LOẠI LỖI
// =====================================================================

// NonRetryableError bọc các lỗi nghiệp vụ không cần retry (VD: Data sai, User không tồn tại)
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string { return e.Err.Error() }
func (e *NonRetryableError) Unwrap() error { return e.Err } // Hỗ trợ errors.Is và errors.As

// Hàm Helper để Dev sử dụng trong Handler
func NewNonRetryableError(err error) error {
	return &NonRetryableError{Err: err}
}

// Hàm kiểm tra xem lỗi này có phải loại cấm Retry không
func isNonRetryable(err error) bool {
	var target *NonRetryableError
	return errors.As(err, &target)
}

// =====================================================================
// 2. CẤU HÌNH RETRY TỐI ƯU
// =====================================================================
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,  // Lần 1 chờ 1s
		MaxBackoff:     15 * time.Second, // Chờ tối đa 15s để không block Kafka quá lâu
	}
}

type DLQPublisher interface {
	PublishToDLQ(ctx context.Context, originalTopic string, event events.IntegrationEvent, errCause error) error
}

// =====================================================================
// 3. HÀM XỬ LÝ CHÍNH ĐẠT CHUẨN (WRAPPER)
// =====================================================================

func HandleWithRetryAndDLQ(
	ctx context.Context,
	event events.IntegrationEvent,
	handler events.EventHandler,
	dlq DLQPublisher,
	cfg RetryConfig,
) error {
	var err error

	// Vòng lặp Retry
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {

		// 1. Thực thi Logic
		err = handler(ctx, event)
		if err == nil {
			return nil // THÀNH CÔNG: Thoát ngay, Kafka sẽ Commit
		}

		// 2. Kiểm tra Server có đang bị tắt không (Graceful Shutdown)
		if ctx.Err() != nil {
			slog.Warn("Tiến trình bị huỷ (Shutdown) lúc đang xử lý", "event_id", event.ID)
			return fmt.Errorf("context canceled processing event %s: %w", event.ID, ctx.Err())
		}

		// 3. Kiểm tra loại lỗi. Nếu là lỗi "Vô phương cứu chữa" -> Nghỉ Retry, ném vào DLQ luôn.
		if isNonRetryable(err) {
			slog.Warn("Lỗi nghiệp vụ không thể phục hồi, bỏ qua Retry",
				"event_id", event.ID, "error", err.Error())
			break
		}

		// 4. Nếu đã là lần thử cuối cùng -> Thoát vòng lặp để xuống bước đẩy DLQ
		if attempt == cfg.MaxRetries {
			break
		}

		// 5. Tính toán thời gian chờ: Exponential Backoff + Jitter
		backoff := calculateBackoff(attempt, cfg.InitialBackoff, cfg.MaxBackoff)

		slog.Warn("Xử lý thất bại, chuẩn bị Retry",
			"event_id", event.ID,
			"attempt", attempt+1,
			"wait_time", backoff.String(),
			"error", err.Error(),
		)

		// 6. CHỜ (SLEEP) AN TOÀN - Lắng nghe tín hiệu Shutdown
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done(): // Có người bấm tắt Server
			timer.Stop()
			return fmt.Errorf("context canceled during backoff wait: %w", ctx.Err())
		case <-timer.C: // Hết thời gian đếm ngược, tiếp tục vòng lặp
			// Đi tiếp lên đầu vòng lặp
		}
	}

	// =====================================================================
	// KẾT THÚC RETRY -> ĐẨY VÀO DLQ
	// =====================================================================
	slog.Error("Event vượt quá số lần Retry HOẶC bị lỗi nghiêm trọng, đẩy vào DLQ",
		"event_id", event.ID,
		"topic", event.Topic,
		"final_error", err.Error(),
	)

	dlqErr := dlq.PublishToDLQ(ctx, event.Topic, event, err)
	if dlqErr != nil {
		slog.Error("CRITICAL: DLQ System Offline. Bắt buộc giữ lại Event",
			"event_id", event.ID, "dlq_error", dlqErr.Error())
		// Trả về lỗi để Hàm Consumer tổng (ở ngoài) KHÔNG COMMIT
		return fmt.Errorf("fatal pipeline breakdown. Handler err: %v. DLQ err: %w", err, dlqErr)
	}

	// DLQ đã nhận hàng an toàn -> Trả về nil để Kafka Commit topic chính, gỡ kẹt cho luồng.
	return nil
}

// =====================================================================
// HÀM HELPER TOÁN HỌC: EXPONENTIAL BACKOFF + JITTER
// =====================================================================
func calculateBackoff(attempt int, initial, max time.Duration) time.Duration {
	// Tăng theo cấp số nhân: 1s -> 2s -> 4s -> 8s
	factor := math.Pow(2, float64(attempt))
	baseWait := time.Duration(float64(initial) * factor)

	if baseWait > max {
		baseWait = max
	}

	// Jitter: Thêm độ nhiễu ngẫu nhiên +/- 20%
	// Tránh trường hợp 1000 message cùng lỗi và cùng Retry ở ĐÚNG 1 thời điểm ở tương lai
	jitter := time.Duration(rand.Float64() * 0.2 * float64(baseWait))
	if rand.Intn(2) == 0 {
		baseWait += jitter
	} else {
		baseWait -= jitter
	}

	return baseWait
}

func HandleBatchWithRetryAndDLQ(
	ctx context.Context,
	batch []events.IntegrationEvent,
	handler events.BatchEventHandler,
	dlq DLQPublisher,
	cfg RetryConfig,
) error {
	var err error

	// 1. Vòng lặp Retry cho toàn bộ Batch
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {

		err = handler(ctx, batch)
		if err == nil {
			return nil // Thành công cả mẻ -> Thoát để Commit
		}

		if ctx.Err() != nil {
			slog.Warn("Context bị huỷ lúc đang xử lý Batch")
			return fmt.Errorf("context canceled processing batch: %w", ctx.Err())
		}

		if isNonRetryable(err) {
			slog.Warn("Lỗi nghiệp vụ không thể phục hồi cho Batch, chuyển hướng sang DLQ", "error", err.Error())
			break
		}

		if attempt == cfg.MaxRetries {
			break
		}

		backoff := calculateBackoff(attempt, cfg.InitialBackoff, cfg.MaxBackoff)
		slog.Warn("Batch xử lý thất bại, chuẩn bị Retry", "attempt", attempt+1, "error", err.Error())

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("context canceled during backoff wait: %w", ctx.Err())
		case <-timer.C:
		}
	}

	// =====================================================================
	// 2. KẾT THÚC RETRY -> ĐẨY TOÀN BỘ BATCH VÀO DLQ
	// =====================================================================
	slog.Error("Batch thất bại hoàn toàn, tiến hành đẩy vào DLQ", "batch_size", len(batch), "error", err.Error())

	var dlqErrs error
	for _, ev := range batch {
		if dlqErr := dlq.PublishToDLQ(ctx, ev.Topic, ev, err); dlqErr != nil {
			// Gom tất cả lỗi DLQ lại (Yêu cầu Go 1.20+)
			dlqErrs = errors.Join(dlqErrs, fmt.Errorf("event %s: %w", ev.ID, dlqErr))
		}
	}

	if dlqErrs != nil {
		slog.Error("CRITICAL: Hệ thống DLQ từ chối nhận Batch", "dlq_errors", dlqErrs.Error())
		return fmt.Errorf("fatal batch pipeline breakdown: %w", dlqErrs)
	}

	// Đã đưa an toàn vào DLQ -> Trả về nil để báo Kafka đi tiếp
	return nil
}
