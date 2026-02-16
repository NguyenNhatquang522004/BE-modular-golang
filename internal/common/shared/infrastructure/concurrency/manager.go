package concurrency

import (
	"context"
	"log"
	"sync"
)

// Service định nghĩa behavior chung cho bất kỳ worker nào (Kafka, Cron, etc)
type Service interface {
	Name() string
	Start(ctx context.Context) error // Hàm này nên là Blocking
	Stop()                           // Hàm này dùng để Graceful Shutdown
}

// Manager quản lý tập trung tất cả Service
type Manager struct {
	services []Service
	wg       sync.WaitGroup
}

// NewManager khởi tạo Manager với danh sách các service
func NewManager(services ...Service) *Manager {
	return &Manager{
		services: services,
	}
}

// StartAll chạy tất cả service trong các goroutine riêng biệt
func (m *Manager) Start(ctx context.Context) {
	for _, s := range m.services {
		m.wg.Add(1)
		// Capture variable cho closure
		service := s

		go func() {
			defer m.wg.Done()
			log.Printf("🚀 Starting Service: %s", service.Name())

			// Start service (Blocking call)
			if err := service.Start(ctx); err != nil {
				log.Printf("❌ Service %s stopped with error: %v", service.Name(), err)
			} else {
				log.Printf("⚠️ Service %s stopped normally", service.Name())
			}
		}()
	}
}

// StopAll gọi lệnh dừng cho từng service và đợi tất cả tắt hẳn
func (m *Manager) Stop() {
	log.Println("🛑 Stopping all background services...")

	// 1. Gửi lệnh Stop cho từng thằng (Async)
	// Lưu ý: Các service nên implement Stop() nhanh gọn hoặc chờ worker của nó
	for _, s := range m.services {
		log.Printf("⏳ Stopping %s...", s.Name())
		s.Stop()
	}

	// 2. Đợi tất cả goroutine trong Start() kết thúc
	m.wg.Wait()
	log.Println("✅ All services stopped gracefully.")
}
