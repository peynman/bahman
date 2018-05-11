package modules

import (
	"avalanche/app/core/app"
	"avalanche/app/core/interfaces"
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
func (_ *SamplePlugin) Initialize() bool {
	return true
}
func (_ *SamplePlugin) Interface() interface{} {
	return nil
}
var PluginInstance interfaces.AvalanchePlugin = new(SamplePlugin)



