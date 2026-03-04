package events

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
)

// Định nghĩa các trạng thái xử lý chuẩn (Enum)

// ConsumerResult là DTO bọc lại kết quả của MỌI handler.
// Ký hiệu [T any] cho phép trường Data chứa bất kỳ kiểu dữ liệu nào (struct, int, string, map...).
type ConsumerResult[T any] struct {
	OriginalEventID string                  `json:"original_event_id"` // ID của event gốc để truy vết (Traceability)
	Topic           string                  `json:"topic"`
	Status          constants.ProcessStatus `json:"status"`

	// Data chứa kết quả ĐẦU RA sau khi xử lý.
	// Ví dụ: Tạo user xong thì trả về UserID, tạo đơn hàng xong trả về Order struct.
	Data T `json:"data,omitempty"`

	// Lưu trữ thông tin lỗi chi tiết (nếu có)
	Error      string `json:"error,omitempty"`
	ErrorCause string `json:"error_cause,omitempty"` // Lỗi gốc từ thư viện/DB

	// Metadata lưu các thông tin phụ trợ cho việc Monitor
	ProcessedAt    time.Time     `json:"processed_at"`
	ProcessingTime time.Duration `json:"processing_time_ms"` // Thời gian chạy task mất bao lâu
	RetryCount     int           `json:"retry_count"`
}

func NewSuccessResult[T any](event IntegrationEvent, data T, duration time.Duration) ConsumerResult[T] {
	return ConsumerResult[T]{
		OriginalEventID: event.ID,
		Topic:           event.Topic,
		Status:          constants.StatusSuccess,
		Data:            data,
		ProcessedAt:     time.Now(),
		ProcessingTime:  duration,
	}
}

// Helper tạo kết quả THẤT BẠI (Dùng để ném vào DLQ)
func NewFailedResult[T any](event IntegrationEvent, err error, retries int, duration time.Duration) ConsumerResult[T] {
	errMessage := "unknown error"
	if err != nil {
		errMessage = err.Error()
	}

	return ConsumerResult[T]{
		OriginalEventID: event.ID,
		Topic:           event.Topic,
		Status:          constants.StatusFailed,
		Error:           errMessage,
		ProcessedAt:     time.Now(),
		ProcessingTime:  duration,
		RetryCount:      retries,
	}
}

// Helper tạo kết quả BỎ QUA (Ví dụ: Đã xử lý rồi, check DB thấy trùng)
func NewSkippedResult[T any](event IntegrationEvent, reason string) ConsumerResult[T] {
	return ConsumerResult[T]{
		OriginalEventID: event.ID,
		Topic:           event.Topic,
		Status:          constants.StatusSkipped,
		Error:           reason, // Ghi lý do bỏ qua vào đây
		ProcessedAt:     time.Now(),
	}
}
