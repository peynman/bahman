package services

type SessionManager interface {
	CreateCookie(name string, data string) error
	Get(name string) error

}