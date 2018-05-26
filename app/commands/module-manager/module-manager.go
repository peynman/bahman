package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"math"
)


type ModulesConsolePlugin struct {
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
func (d *ModulesConsole) ShowModulesList(console interfaces.ConsoleApp, avModules [] interfaces.Module, callback func(module interfaces.Module)) {
	var items = make([]interfaces.ListItem, int(math.Max(float64(len(avModules)), 1)))

	if len(avModules) > 0 {
		for index, module := range avModules {
			caller := func(lModule interfaces.Module) interfaces.ListItem {
				return interfaces.ListItem{
					Title: lModule.Title(),
					Description: lModule.Description(),
					Shortcut: 0,
					Callback: func() {
						callback(lModule)
					},
				}
			}
			items[index] = caller(module)
		}
	} else {
		items[0] = interfaces.ListItem{
			Title: "No module found in this list!",
			Description:"",
			Shortcut: 0,
			Callback: func() {
				console.Back()
			},
		}
	}
	console.SetPage("modules_list", console.MakeList("Available modules", items), true)
}
func (d *ModulesConsole) OnSelected(console interfaces.ConsoleApp) {
	console.SetPage(
		"modules_manager_main",
		console.MakeList("Avalanche Module Management", []interfaces.ListItem {
			{
				"List",
				"List available module-manager",
				'l',
				func() {
					avModules := d.services.Modules().List()
					d.ShowModulesList(console, avModules, func(module interfaces.Module) {
						console.SetPage("module_about", console.MakeModal(interfaces.ModalWindow{
							Content:module.Title() + " - " + module.Description() + " - Version: " + module.Version(),
							Buttons:[]string {"Done"},
							Callback: func(buttonIndex int) {
								console.Back()
							},
						}), false)
					})
				},
			},
			{
				"Install",
				"Install a available",
				'i',
				func() {
					avModules := d.services.Modules().NotInstalled()
					d.ShowModulesList(console, avModules, func(module interfaces.Module) {
						console.Ask("Are you sure you want to activate this module?", func(yes bool) {
							if yes {
								err := d.services.Modules().Install(module)
								if err != nil {

								}
							}
							console.Back()
						})
					})
				},
			},
			{
				"Activate",
				"Activate a module",
				'a',
				func() {
					avModules := d.services.Modules().Deactivated()
					d.ShowModulesList(console, avModules, func(module interfaces.Module) {

					})
				},
			},
			{
				"Deactivate",
				"Deactivate a module; but don not remove its files or data",
				'd',
				func() {
					avModules := d.services.Modules().Activated()
					d.ShowModulesList(console, avModules, func(module interfaces.Module) {

					})
				},
			},
			{
				"Purge",
				"Remove all module files and rollback all migrations",
				'p',
				func() {
					avModules := d.services.Modules().Installed()
					d.ShowModulesList(console, avModules, func(module interfaces.Module) {

					})
				},
			},
		}),
		true)
}



