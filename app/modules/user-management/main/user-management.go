package main

import (
	"avalanche/app/core/app"
	"avalanche/app/core/interfaces"
	"avalanche/app/modules/user-management"
)

type UserManagementPlugin struct {
}
func (_ *UserManagementPlugin) Version() string {
	return "1.0.0"
}
func (_ *UserManagementPlugin) VersionCode() int {
	return 1
}
func (_ *UserManagementPlugin) AvalancheVersionCode() int {
	return app.VersionCode
}
func (_ *UserManagementPlugin) Title() string {
	return "User Management"
}
func (_ *UserManagementPlugin) Description() string {
	return "Default User-Management system"
}
func (_ *UserManagementPlugin) Initialize() bool {
	userManagementModule = new(UserManagementModule)
	return true
}
func (_ *UserManagementPlugin) Interface() interface{} {
	return userManagementModule
}
var PluginInstance interfaces.AvalanchePlugin = new(UserManagementPlugin)

type UserManagementModule struct {
}
var userManagementModule *UserManagementModule
var _ interfaces.Module = (*UserManagementModule)(nil)

func (_ *UserManagementModule) Migrations() []interfaces.Migratable {
	return []interfaces.Migratable {
		new(user_management.UsersTable),
	}
}