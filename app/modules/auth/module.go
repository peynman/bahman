package auth

import "github.com/peyman-abdi/avalanche/app/interfaces/services"

type AuthenticationModule struct {
	Services services.Services
}

func (*AuthenticationModule) Title() string { return "Authentication Modules" }

func (*AuthenticationModule) Description() string { return "Default user/role/permission" }

func (*AuthenticationModule) Version() string { return "1.0.0" }

func (*AuthenticationModule) Routes() []*services.Route {
	return nil
}

func (*AuthenticationModule) GroupsHandlers() []*services.RouteGroup {
	return nil
}

func (*AuthenticationModule) Templates() []*services.Template {
	return nil
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

var _ services.Module = (*AuthenticationModule)(nil)
