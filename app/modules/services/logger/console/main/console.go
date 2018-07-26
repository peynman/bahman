package main

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/logger/console"
)

type bahmanPluginConsole struct {
	instance services.Services
	console *console.LogChannel
}

func (a *bahmanPluginConsole) Initialize(services services.Services) bool {
	a.instance = services
	a.console = console.New(a.instance)
	return true
}
func (a *bahmanPluginConsole) Interface() interface{} {
	return a.console
}

var PluginInstance services.bahmanPlugin = new(bahmanPluginConsole)
