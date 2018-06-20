package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/peyman-abdi/avalanche/app/modules/auth"
)

type AuthManagerPlugin struct {
	module *auth.AuthenticationModule
}

func (a *AuthManagerPlugin) Initialize(services core.Services) bool {
	a.module = new(auth.AuthenticationModule)
	a.module.Services = services
	return true
}
func (a *AuthManagerPlugin) Interface() interface{} {
	return a.module
}

var PluginInstance core.AvalanchePlugin = new(AuthManagerPlugin)
