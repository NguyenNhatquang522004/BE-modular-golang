package IConsumer

type IWorkerBlock interface {
	CreateBlock() error
	DeleteBlock() error
	UpdateBlock() error
}