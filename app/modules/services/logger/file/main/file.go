package main

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/logger/file"
)
type bahmanPluginFile struct {
	instance services.Services
	fileLogger *file.LogChannel
}

func (a *bahmanPluginFile) Initialize(services services.Services) bool {
	a.instance = services
	a.fileLogger = file.New(a.instance)
	return true
}
func (a *bahmanPluginFile) Interface() interface{} {
	return a.fileLogger
}

var PluginInstance services.bahmanPlugin = new(bahmanPluginFile)
