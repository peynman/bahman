package routes

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/auth/middlewares"
	"github.com/peyman-abdi/bahman/app/interfaces/auth"
	"errors"
)

const (
	NameRouteViewSignIn = "sign-in-view"
	NameRouteAPISignIn = "sign-in-api"
)
func attempt(a services.Services, request services.Request) (auth.User,error) {
	if authier, ok := a.GetByName("auth").(auth.Authenticator); ok {
		user, err := authier.Attempt(request.Session(), request.GetBodyAsMap())
		if err != nil {
			return nil, err
		}
		authier.Remember(request.Session(), user)

		return user, err
	}
	return nil, errors.New("auth module not found")
}
func NewSignInViewRoute(a services.Services) *services.Route {
	return &services.Route{
		Name:  NameRouteViewSignIn,
		Group: a.Config().GetString("auth.routes.views.group", ""),
		MiddleWares: append([]string {
			middlewares.NameMiddleWareGuest,
		}, a.Config().GetStringArray("auth.routes.views.middleWares", []string{})...),
		Methods: services.POST|services.GET,
		Urls: []string {"sign-in"},
		Handle: func(request services.Request, response services.Response) error {
			if request.GetMethod() & services.GET != 0 {
				return response.View(auth.NameTemplateSignIn, nil)
			} else if request.GetMethod() & services.POST != 0 {
				_, err := attempt(a, request)
				if err != nil {
					if appErr, ok := err.(services.ApplicationError); ok {
						if appErr.IsErrorType(a.GetByName(auth.NameService)) {
							switch appErr.Code() {
							case auth.ErrInvalidPassword,
								 auth.ErrUserNotFound:
								 	return response.View(auth.NameTemplateSignIn, map[string]interface{} {
								 		"error": a.Localization().L("auth.errors.user_or_pass_wrong"),
									})
							}
						}
					}

					return response.View(auth.NameTemplateSignIn, map[string]interface{} {
						"error": a.Localization().L("auth.errors.general"),
						"err": err,
					})
				}

				return response.Redirect(a.Router().Url(a.Config().GetString("auth.routes.home", "")), 31)
			}
			return errors.New("method not allowed")
		},
	}
}

func NewSignInApiRoute(a services.Services) *services.Route {
	return &services.Route{
		Name:  NameRouteAPISignIn,
		Group: a.Config().GetString("auth.routes.api.group", "api"),
		MiddleWares: append([]string {
			middlewares.NameMiddleWareGuest,
		}, a.Config().GetStringArray("auth.routes.api.middleWares", []string{})...),
		Methods: services.POST,
		Handle: func(request services.Request, response services.Response) error {
			user, err := attempt(a, request)
			if err != nil {

			}
			return err
		},
	}
}
