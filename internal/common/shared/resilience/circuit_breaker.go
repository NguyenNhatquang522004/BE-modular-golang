package resilience

import (
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/sony/gobreaker"
)

// BreakerProvider interface để inject
type BreakerProvider interface {
	New(name string, cfg *configs.Config) *gobreaker.CircuitBreaker
}

type breakerProviderImpl struct{}

func NewBreakerProvider() BreakerProvider {
	return &breakerProviderImpl{}
}

func (p *breakerProviderImpl) New(name string, cfg *configs.Config) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    time.Duration(cfg.CircuitBreaker.Interval) * time.Millisecond,
		Timeout:     time.Duration(cfg.CircuitBreaker.RecoveryTimeoutSeconds) * time.Second,
		// ReadyToTrip: Logic quyết định khi nào thì ngắt mạch (Open)
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			// Ngắt mạch nếu thất bại > 60% và có ít nhất 10 request
			return counts.Requests >= cfg.CircuitBreaker.MaxRequests && failureRatio >= cfg.CircuitBreaker.FailureThreshold/100.0
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Tại đây nên bắn log hoặc metric (Prometheus)
			log.Printf("Circuit Breaker '%s' changed from %s to %s", name, from, to)
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
