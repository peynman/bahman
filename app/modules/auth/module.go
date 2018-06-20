package auth

import "github.com/peyman-abdi/avalanche/app/interfaces/core"

type AuthenticationModule struct {
	Services core.Services
}

func (*AuthenticationModule) Title() string { return "Authentication Modules" }

func (*AuthenticationModule) Description() string { return "Default user/role/permission" }

func (*AuthenticationModule) Version() string { return "1.0.0" }

func (*AuthenticationModule) Routes() []*core.Route {
	return nil
}

func (*AuthenticationModule) MiddleWares() []*core.MiddleWare {
	return nil
}

func (*AuthenticationModule) GroupsHandlers() []*core.RouteGroup {
	return nil
}

func (*AuthenticationModule) Templates() []*core.Template {
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

var _ core.Module = (*AuthenticationModule)(nil)
