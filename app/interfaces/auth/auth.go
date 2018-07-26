package auth

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
)

const (
	ErrUserNotFound = 1
	ErrInvalidPassword = 2

	NameService = "auth"
	NameTemplateSignIn = "sign-in"
	NameTemplateSignUp = "sign-up"
	NameTemplateResetPassword = "reset-pass"
	NameTemplateUpdatePassword = "update-pass"
)
type Authenticator interface {
	IsValid(session services.Session) bool
	Guest(session services.Session) bool
	User(session services.Session) (User, error)

	Attempt(session services.Session, credentials map[string]string) (User, error)
	Remember(session services.Session, user User)
	Forget(session services.Session)

	Use(provider UserProvider)
}
type User interface {
	UID() string
}
type UserProvider interface {
	FindByID(id string) (User, error)
	FindByCredentials(credential map[string]string) (User, error)
	IsValidCredentials(user User, credential map[string]string) bool
}
