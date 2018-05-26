package modules

import (
	"github.com/peyman-abdi/avalanche/app/core/app"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type SamplePlugin struct {
}
func (_ *SamplePlugin) Version() string {
	return "1.0.0"
}
func (_ *SamplePlugin) VersionCode() int {
	return 1
}
func (_ *SamplePlugin) AvalancheVersionCode() int {
	return app.VersionCode
}
func (_ *SamplePlugin) Title() string {
	return "File Logger"
}
func (_ *SamplePlugin) Description() string {
	return "File driver for Avalanche logger"
}
func (_ *SamplePlugin) Initialize(services interfaces.Services) bool {
	return true
}
func (_ *SamplePlugin) Interface() interface{} {
	return nil
}
var PluginInstance interfaces.AvalanchePlugin = new(SamplePlugin)



