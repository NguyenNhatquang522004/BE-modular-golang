package IRepositoryShare

type IEmail interface {
	SendEmail(to string, subject string, body string) error
}
