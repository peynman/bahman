package auth

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/auth/migrations"
	"github.com/peyman-abdi/bahman/app/modules/auth/middlewares"
	"github.com/peyman-abdi/bahman/app/interfaces/auth"
)

type AuthenticationModule struct {
	Instance services.Services
}

func (*AuthenticationModule) Title() string { return "Authentication Modules" }

func (*AuthenticationModule) Description() string { return "Default user/role/permission" }

func (*AuthenticationModule) Version() string { return "1.0.0" }

func (*AuthenticationModule) Migrations() []services.Migratable {
	return []services.Migratable{
		new(migrations.AuthUserMigrate),
	}
}

func (a *AuthenticationModule) MiddleWares() []*services.MiddleWare {
	return []*services.MiddleWare {
		middlewares.NewAuthMiddleWare(a.Instance),
		middlewares.NewGuestMiddleWare(a.Instance),
	}
}
func (*AuthenticationModule) Routes() []*services.Route {
	return nil
}

func (*AuthenticationModule) GroupsHandlers() []*services.RouteGroup {
	return nil
}

func (a *AuthenticationModule) Services() map[string]func()interface{} {
	return map[string]func() interface {} {
		auth.NameService: func() interface {} {
			return NewAuthenticator(a.Instance)
		},
	}
}

func (*AuthenticationModule) Activated() bool {
	return true
}

func (*AuthenticationModule) Installed() bool {
	return true
}

func (*AuthenticationModule) Deactivated() {
}

func (*AuthenticationModule) Purged() {
}

func (*AuthenticationModule) Templates() []*services.Template {
	return []*services.Template {
		{
			Name: auth.NameTemplateSignIn,
			Path: "auth/sign-in.jet",
		},
		{
			Name: auth.NameTemplateSignUp,
			Path: "auth/sign-up.jet",
		},
		{
			Name: auth.NameTemplateResetPassword,
			Path: "auth/reset-pass.jet",
		},
		{
			Name: auth.NameTemplateUpdatePassword,
			Path: "auth/update-pass.jet",
		},
	}
}

var _ services.Module = (*AuthenticationModule)(nil)

