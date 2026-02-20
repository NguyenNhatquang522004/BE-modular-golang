package concurrency

import (
	"context"
	"sync"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"golang.org/x/sync/semaphore"
)

type WorkerPool struct {
	sem            *semaphore.Weighted
	wg             sync.WaitGroup
	maxConcurrency int64
	cfg            *configs.Config
}

func NewWorkerPool(cfg *configs.Config) *WorkerPool {
	maxConcurrency := int64(cfg.WorkerPool.WORKER_POOL_SIZE_MAX)
	return &WorkerPool{
		sem:            semaphore.NewWeighted(maxConcurrency),
		maxConcurrency: maxConcurrency,
		cfg:            cfg,
	}
}

// Run: Chạy task async nhưng có kiểm soát số lượng và tracking vòng đời
func (p *WorkerPool) Run(ctx context.Context, task func()) error {
	// 1. Xin slot (Block nếu đã full slot)
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return err
	}

	p.wg.Add(1) // 2. Đánh dấu có 1 task đang chạy

	go func() {
		defer p.sem.Release(1) // 4. Trả slot
		defer p.wg.Done()      // 5. Báo cáo xong task

		// Recover panic để an toàn
		defer func() {
			if r := recover(); r != nil { /* Log error */
			}
		}()

		task() // 3. Thực thi logic
	}()

	return nil
}

// Wait: Dùng khi Shutdown server
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
