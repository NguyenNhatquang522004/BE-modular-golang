// pkg/scheduler/scheduler.go
package IRepositoryShare

// JobFunc là định dạng của hàm sẽ được thực thi khi đến giờ
type JobFunc func()

// Scheduler định nghĩa các hành động mà bộ lập lịch có thể làm
type IScheduler interface {
	// Cho phép đăng ký một hàm với thời gian chạy cụ thể
	ScheduleJob(cronExpr string, job JobFunc) error
	Start()
	Stop()
}
