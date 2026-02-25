// internal/infrastructure/cron/cron_scheduler.go
package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron *cron.Cron
}

// NewCronScheduler khởi tạo bộ lập lịch với múi giờ Local
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(cron.WithLocation(time.Local)),
	}
}

func (s *CronScheduler) ScheduleJob(cronExpr string, job func()) error {
	_, err := s.cron.AddFunc(cronExpr, job)
	return err
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
}
