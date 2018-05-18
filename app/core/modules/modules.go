package modules

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"reflect"
)

type moduleManagerImpl struct {
	AvailableModules []interfaces.Module
	ModuleTableName string
	services interfaces.Services
}

func Initialize(config interfaces.Config) interfaces.ModuleManager {
	module := new(moduleManagerImpl)

	module.ModuleTableName = config.GetString("modules.table", "modules")

	return module
}

func (m *moduleManagerImpl) LoadModules(services interfaces.Services) {
	modules := services.App().InitAvalanchePlugins(services.App().ModulesPath("app"), services)

	m.AvailableModules = make([]interfaces.Module, len(modules))
	for index, moduleInterface := range modules {
		m.AvailableModules[index] = moduleInterface.Interface().(interfaces.Module)
	}
}

func (m *moduleManagerImpl) IsActive(module interfaces.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)
	if model != nil {
		return model.Flags & INSTALLED != 0 && model.Flags & ACTIVE != 0
	}

	return false
}
func (m *moduleManagerImpl) IsInstalled(module interfaces.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)
	if model != nil {
		return model.Flags & INSTALLED != 0
	}

	return false
}
func (m *moduleManagerImpl) List() []interfaces.Module {
	return m.AvailableModules
}

func (m *moduleManagerImpl) Install(module interfaces.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)

	m.services.Migrator().Migrate(module.Migrations())

	if model != nil {
		if model.Flags & INSTALLED != 0 {
			return ModuleError(module, "Target module is already installed")
		}

		model.Flags |= INSTALLED
		m.services.Repository().UpdateEntity(model)
	} else {
		model = &ModuleModel{
			Flags: INSTALLED,
			Interface: reflect.TypeOf(module).String(),
		}
		err := m.services.Repository().Insert(model)
		if err != nil {

		}
	}

	if !module.Installed() {
		m.Purge(module)
		return ModuleError(module, "Module failed to install")
	}

	return nil
}
func (m *moduleManagerImpl) Activate(module interfaces.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)

	m.services.Migrator().Migrate(module.Migrations())

	if model == nil || model.Flags & INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be activated")
	}
	if model.Flags & ACTIVE != 0 {
		return ModuleError(module, "Target module is already active")
	}

	model.Flags |= ACTIVE
	m.services.Repository().UpdateEntity(model)

	if !module.Activated() {
		m.Deactivate(module)

		model.Flags |= ACTIVE
		m.services.Repository().UpdateEntity(model)
	}

	return nil
}
func (m *moduleManagerImpl) Purge(module interfaces.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)

	m.services.Migrator().Migrate(module.Migrations())

	if model == nil || model.Flags & INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be deactivated")
	}

	model.Flags &= ^(ACTIVE | INSTALLED)
	m.services.Repository().UpdateEntity(model)

	module.Purged()

	return nil
}
func (m *moduleManagerImpl) Deactivate(module interfaces.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module)).Get(model)

	m.services.Migrator().Migrate(module.Migrations())

	if model == nil || model.Flags & INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be deactivated")
	}
	if model.Flags & ACTIVE == 0 {
		return ModuleError(module,"Target module is not activated yet and can not be deactivated")
	}

	model.Flags &= ^ACTIVE
	m.services.Repository().UpdateEntity(model)

	module.Deactivated()

	return nil
}


