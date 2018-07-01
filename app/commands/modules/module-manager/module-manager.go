package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"math"
	"fmt"
)

type ModulesConsolePlugin struct {
}

func (_ *ModulesConsolePlugin) Initialize(services services.Services) bool {
	modulesConsole = new(ModulesConsole)
	modulesConsole.services = services
	return true
}
func (_ *ModulesConsolePlugin) Interface() interface{} {
	return modulesConsole
}

var PluginInstance services.AvalanchePlugin = new(ModulesConsolePlugin)

type ModulesConsole struct {
	services services.Services
}

var _ services.ConsolePage = (*ModulesConsole)(nil)
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
func (d *ModulesConsole) ShowModulesList(console services.ConsoleApp, avModules []services.Module, callback func(module services.Module)) {
	var items = make([]services.ListItem, int(math.Max(float64(len(avModules)), 1)))

	if len(avModules) > 0 {
		for index, module := range avModules {
			caller := func(lModule services.Module) services.ListItem {
				return services.ListItem{
					Title:       lModule.Title(),
					Description: lModule.Description(),
					Shortcut:    0,
					Callback: func() {
						callback(lModule)
					},
				}
			}
			items[index] = caller(module)
		}
	} else {
		items[0] = services.ListItem{
			Title:       "No module found in this list!",
			Description: "",
			Shortcut:    0,
			Callback: func() {
				console.Back()
			},
		}
	}
	console.SetPage("modules_list", console.MakeList("Available modules", items), true)
}
func (d *ModulesConsole) OnSelected(console services.ConsoleApp) {
	console.SetPage(
		"modules_manager_main",
		console.MakeList("Avalanche Module Management", []services.ListItem{
			{
				"List",
				"List available module-manager",
				'l',
				func() {
					avModules := d.services.Modules().List()
					d.ShowModulesList(console, avModules, func(module services.Module) {
						console.SetPage("module_about", console.MakeModal(services.ModalWindow{
							Content: module.Title() + " - " + module.Description() + " - Version: " + module.Version(),
							Buttons: []string{"Done"},
							Callback: func(buttonIndex int) {
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
					d.ShowModulesList(console, avModules, func(module services.Module) {
						console.Ask("Are you sure you want to install this module?", func(yes bool) {
							if yes {
								err := d.services.Modules().Install(module)
								if err != nil {
									console.Error(fmt.Sprintf("Could not install module: %s", err), func() {
									})
								} else {
									console.Success(fmt.Sprintf("Installation was successful"), func() {
									})
								}
							}
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
					d.ShowModulesList(console, avModules, func(module services.Module) {
						console.Ask("Are you sure you want to activate this module?", func(yes bool) {
							if yes {
								err := d.services.Modules().Activate(module)
								if err != nil {
									console.Error(fmt.Sprintf("Could not activate module: %s", err), func() {
									})
								} else {
									console.Success(fmt.Sprintf("Activation was successful"), func() {
									})
								}
							}
						})
					})
				},
			},
			{
				"Deactivate",
				"Deactivate a module; but don not remove its files or data",
				'd',
				func() {
					avModules := d.services.Modules().Activated()
					d.ShowModulesList(console, avModules, func(module services.Module) {
						console.Ask("Are you sure you want to deactivate this module?", func(yes bool) {
							if yes {
								err := d.services.Modules().Deactivate(module)
								if err != nil {
									console.Error(fmt.Sprintf("Could not deactivate module: %s", err), func() {
									})
								} else {
									console.Success(fmt.Sprintf("Deactivation was successful"), func() {
									})
								}
							}
						})
					})
				},
			},
			{
				"Purge",
				"Remove all module files and rollback all migrations",
				'p',
				func() {
					avModules := d.services.Modules().Installed()
					d.ShowModulesList(console, avModules, func(module services.Module) {
						console.Ask("Are you sure you want to purge this module? all stored changes will be lost.", func(yes bool) {
							if yes {
								err := d.services.Modules().Deactivate(module)
								if err != nil {
									console.Error(fmt.Sprintf("Could not purge module: %s", err), func() {
									})
								} else {
									console.Success(fmt.Sprintf("Purge was successful"), func() {
									})
								}
							}
						})
					})
				},
			},
		}),
		true)
}
