package IConsumer

type IWorkerFollow interface {
	CreateFollow() error
	DeleteFollow() error
	UpdateFollow() error
}
