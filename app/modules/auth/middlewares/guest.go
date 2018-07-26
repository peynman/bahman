package middlewares

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/interfaces/auth"
)

const (
	NameMiddleWareGuest = "guest"
)
func NewGuestMiddleWare(a services.Services) *services.MiddleWare {
	return &services.MiddleWare{
		Name: NameMiddleWareGuest,
		Priority: 10,
		Handler: func(request services.Request, response services.Response) error {
			if auth, ok := a.GetByName("auth").(auth.Authenticator); ok {
				user, err := auth.User(request.Session())
				if err == nil && user != nil {
					request.SetValue("auth.id", user.UID())
					request.SetValue("auth.user", user)
				}
			}
			return a.App().Error(a, auth.ErrUserNotFound,"")
		},
	}
}

