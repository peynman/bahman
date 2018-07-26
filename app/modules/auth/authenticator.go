package auth

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/interfaces/auth"
	"github.com/uniplaces/carbon"
	"time"
)

const DefaultExpiration = 60 * 24 * 60 * 60

type authenticatorImpl struct {
	provider auth.UserProvider
	cookieName string
	instance services.Services
}
func (a *authenticatorImpl) IsValid(session services.Session) bool {
	params := session.GetAsMap("auth")
	if params == nil {
		return false
	}
	if expiresString, ok := params["expires"].(string); ok {
		expires, err := carbon.CreateFromFormat("Y-m-d H:i:s", expiresString, a.instance.Config().GetAsString("app.timezone", "UTC"))
		if err != nil {
			a.instance.Logger().ErrorFields("Invalid auth session expiration format", map[string]interface{} {
				"err": err,
				"value": expiresString,
			})
			return false
		}
		if carbon.Now().DiffInSeconds(expires, false) < 0 {
			return false
		}
	}
	return true
}

func (*authenticatorImpl) Guest(session services.Session) bool {
	panic("implement me")
}

func (*authenticatorImpl) User(session services.Session) (auth.User, error) {
	panic("implement me")
}

func (a *authenticatorImpl) Remember(session services.Session, user auth.User) {
	session.Set("auth", map[string]interface{} {
		"expires": carbon.Now().Add(time.Duration(a.instance.Config().GetInt("auth.session.expires", DefaultExpiration))),
		"id": user.UID(),
	})
}

func (a *authenticatorImpl) Attempt(session services.Session, credentials map[string]string) (auth.User, error) {
	return a.AttemptWithProvider(session, credentials, a.provider)
}

func (a *authenticatorImpl) AttemptWithProvider(session services.Session, credentials map[string]string, provider auth.UserProvider) (auth.User, error) {
	user, err := provider.FindByCredentials(credentials)
	if err != nil {
		return nil, err
	}
	if provider.IsValidCredentials(user, credentials) {
		return user, nil
	}

	return nil, a.instance.App().Error(a, auth.ErrInvalidPassword, a.instance.Localization().L("auth.errors.wrong_password"))
}

func (*authenticatorImpl) Forget(session services.Session) {
	session.Delete("auth")
}

func (a *authenticatorImpl) Use(provider auth.UserProvider) {
	a.provider = provider
}

func NewAuthenticator(services services.Services) auth.Authenticator {
	a := new(authenticatorImpl)
	a.instance = services
	a.provider = NewAuthUserProvider(services)
	return a
}

