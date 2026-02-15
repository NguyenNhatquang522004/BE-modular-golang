package IWorker

type IWorkerBlock interface {
	CreateBlock() error
	DeleteBlock() error
	UpdateBlock() error
}