package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"reflect"
	"fmt"
)

type ModulesConsolePlugin struct {
}
func (_ *ModulesConsolePlugin) Version() string {
	return "1.0.0"
}
func (_ *ModulesConsolePlugin) VersionCode() int {
	return 1
}
func (_ *ModulesConsolePlugin) AvalancheVersionCode() int {
	return 1
}
func (_ *ModulesConsolePlugin) Title() string {
	return "Module Manager"
}
func (_ *ModulesConsolePlugin) Description() string {
	return "Modules management tools for CLI"
}
func (_ *ModulesConsolePlugin) Initialize(services interfaces.Services) bool {
	modulesConsole = new(ModulesConsole)
	modulesConsole.services = services
	return true
}
func (_ *ModulesConsolePlugin) Interface() interface{} {
	return modulesConsole
}
var PluginInstance interfaces.AvalanchePlugin = new(ModulesConsolePlugin)

type ModulesConsole struct {
	services interfaces.Services
}
var _ interfaces.ConsolePage = (*ModulesConsole)(nil)
var modulesConsole *ModulesConsole

func (d *ModulesConsole) Title() string {
	return "Modules"
}
func (d *ModulesConsole) Description() string {
	return "Manage modules from CLI"
}
func (d *ModulesConsole) Priority() int {
	return 10
}
func (d *ModulesConsole) OnSelected(console interfaces.ConsoleApp) {
	console.SetPage(
		console.MakeList("Avalanche Module Management", []interfaces.ListItem {
			{
				"List",
				"List available module-manager",
				'l',
				func() {
					avModules := d.services.Modules().List()
					var items = make([]interfaces.ListItem, len(avModules))
					for index, module := range avModules {
						items[index] = interfaces.ListItem{
							Title: reflect.TypeOf(module).String(),
							Description: "",
							Shortcut: 0,
							Callback: func() {

							},
						}
					}
					console.SetPage(console.MakeList("Available modules", items), true)
				},
			},
			{
				"Install",
				"Install a available",
				'i',
				func() {
					fmt.Println("item clicked inside 1")
				},
			},
			{
				"Activate",
				"Activate a module",
				'a',
				func() {
				},
			},
			{
				"Deactivate",
				"Deactivate a module; but don not remove its files or data",
				'd',
				func() {
				},
			},
			{
				"Purge",
				"Remove all module files and rollback all migrations",
				'p',
				func() {
				},
			},
		}),
		true)
}



