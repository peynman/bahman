package main

import (
	"github.com/rivo/tview"
	"avalanche/app/core/interfaces"
	"avalanche/app/core/app"
)

type DatabaseConsolePlugin struct {
}
func (_ *DatabaseConsolePlugin) Version() string {
	return "1.0.0"
}
func (_ *DatabaseConsolePlugin) VersionCode() int {
	return 1
}
func (_ *DatabaseConsolePlugin) AvalancheVersionCode() int {
	return app.VersionCode
}
func (_ *DatabaseConsolePlugin) Title() string {
	return "File Logger"
}
func (_ *DatabaseConsolePlugin) Description() string {
	return "File driver for Avalanche logger"
}
func (_ *DatabaseConsolePlugin) Initialize() bool {
	databaseConsole = new(DatabaseConsole)
	return true
}
func (_ *DatabaseConsolePlugin) Interface() interface{} {
	return databaseConsole
}
var PluginInstance interfaces.AvalanchePlugin = new(DatabaseConsolePlugin)

type DatabaseConsole struct {
}
var _ interfaces.ConsolePage = (*DatabaseConsole)(nil)
var databaseConsole *DatabaseConsole

func (d *DatabaseConsole) Title() string {
	return "Database"
}
func (d *DatabaseConsole) Description() string {
	return "Migrations & Seeding"
}
func (d *DatabaseConsole) Priority() int {
	return 50
}
func (d *DatabaseConsole) OnSelected(console interfaces.ConsoleApp) {
	commandsList := tview.NewList()
	commandsList.AddItem("Create Migration", "Create new migration", 'c', func() {

	})
	commandsList.AddItem("Migrate", "Commit existing migrations", 'm', func() {

	})
	commandsList.AddItem("Rollback", "Rollback last commit", 'r', func() {

	})
	commandsList.AddItem("Seed Database", "Run database seeds", 's', func() {

	})
	commandsList.AddItem("Back", "Back to main menu", 'b', func() {
		console.Back()
	})
	commandsList.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Database Migrations & Seeding")
	console.SetPage(commandsList, true)
}


