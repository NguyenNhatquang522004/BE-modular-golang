package IRepositoryShare

import "context"

type IWorkerPool interface {
	Run(ctx context.Context, task func()) error
	Wait()
}
