package main

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/peyman-abdi/avalanche/app/modules/user-management"
)

type UserManagementPlugin struct {
}
func (_ *UserManagementPlugin) Initialize(services interfaces.Services) bool {
	userManagementModule = new(UserManagementModule)
	userManagementModule.services = services
	return true
}
func (_ *UserManagementPlugin) Interface() interface{} {
	return userManagementModule
}
var PluginInstance interfaces.AvalanchePlugin = new(UserManagementPlugin)

type UserManagementModule struct {
	services interfaces.Services
}
var userManagementModule *UserManagementModule
var _ interfaces.Module = (*UserManagementModule)(nil)

func (u *UserManagementModule) Title() string {
	return "Avalanche User Management"
}

func (u *UserManagementModule) Description() string {
	return "Minimal users/roles/permissions system for avalanche"
}
func (u *UserManagementModule) Version() string {
	return "1.0.0"
}

func (_ *UserManagementModule) Migrations() []interfaces.Migratable {
	return []interfaces.Migratable {
		new(user_management.UsersTable),
	}
}
func (u *UserManagementModule) Installed() bool {
	return true
}
func (u *UserManagementModule) Activated() bool {
	return true
}
func (u *UserManagementModule) Deactivated() {
}
func (u *UserManagementModule) Purged() {
}