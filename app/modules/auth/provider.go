package auth

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/interfaces/auth"
	"github.com/peyman-abdi/bahman/app/modules/auth/models"
)

type authUserProvider struct {
	instance services.Services
}

func (a *authUserProvider) FindByID(id string) (auth.User, error) {
	var user *models.UserAuth
	err := a.instance.Repository().Query(&models.UserAuth{}).Where("id = ?", id).GetFirst(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, a.instance.App().Error(a, auth.ErrUserNotFound, a.instance.Localization().L("auth.errors.user_not_found"))
	}
	return user, nil
}

func (a *authUserProvider) FindByCredentials(credential map[string]string) (auth.User, error) {
	var user *models.UserAuth
	if username, ok := credential["username"]; ok {
		err := a.instance.Repository().Query(&models.UserAuth{}).Where("username = ?", username).GetFirst(&username)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (a *authUserProvider) IsValidCredentials(user auth.User, credential map[string]string) bool {
	if authUser, ok := user.(*models.UserAuth); ok {
		if password, ok := credential["password"]; ok {
			return a.instance.Hash().Compare(password+authUser.Salt, authUser.Password)
		}
	}
	return false
}

func NewAuthUserProvider(services services.Services) auth.UserProvider {
	return &authUserProvider{
		instance: services,
	}
}



