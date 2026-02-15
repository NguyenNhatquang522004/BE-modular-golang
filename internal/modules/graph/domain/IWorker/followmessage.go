package IWorker

type IWorkerFollow interface {
	CreateFollow() error
	DeleteFollow() error
	UpdateFollow() error
}
