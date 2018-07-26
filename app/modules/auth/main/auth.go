package main

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/auth"
)

type AuthManagerPlugin struct {
	module *auth.AuthenticationModule
}

func (a *AuthManagerPlugin) Initialize(services services.Services) bool {
	a.module = new(auth.AuthenticationModule)
	a.module.Instance = services
	return true
}
func (a *AuthManagerPlugin) Interface() interface{} {
	return a.module
}

//PluginInstance
var PluginInstance services.bahmanPlugin = new(AuthManagerPlugin)
