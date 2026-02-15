package IWorker

type IWorkerFriend interface {
	CreateFriendship() error
	DeleteFriendship() error
	UpdateFriendship() error
}
