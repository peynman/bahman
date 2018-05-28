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

func Initialize(config interfaces.Config, migrator interfaces.Migrator) interfaces.ModuleManager {
	module := new(moduleManagerImpl)

	module.ModuleTableName = config.GetString("modules.table", "modules")

	if !migrator.HasTable(&ModuleModel{}) {
		migrator.AutoMigrate(&ModuleModel{})
	}

	return module
}

func (m *moduleManagerImpl) LoadModules(services interfaces.Services) {
	modules := services.App().InitAvalanchePlugins(services.App().ModulesPath("app"), services)

	m.services = services
	m.AvailableModules = make([]interfaces.Module, len(modules))
	for index, moduleInterface := range modules {
		m.AvailableModules[index] = moduleInterface.Interface().(interfaces.Module)
	}
}

func (m *moduleManagerImpl) IsActive(module interfaces.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if model != nil {
		return model.Flags & INSTALLED != 0 && model.Flags & ACTIVE != 0
	}

	return false
}
func (m *moduleManagerImpl) IsInstalled(module interfaces.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if model != nil {
		return model.Flags & INSTALLED != 0
	}

	return false
}
func (m *moduleManagerImpl) List() []interfaces.Module {
	return m.AvailableModules
}
func (m *moduleManagerImpl) Activated() []interfaces.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("(flags & ?) != 0", ACTIVE).Get(&models)

	var modulesList []interfaces.Module
	for _, model := range models {
		for _, module := range m.AvailableModules {
			if reflect.TypeOf(module).String() == model.Interface {
				modulesList = append(modulesList, module)
				break
			}
		}
	}

	return modulesList
}
func (m *moduleManagerImpl) Installed() []interfaces.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("(flags & ?) != 0", INSTALLED).Get(&models)

	var modulesList []interfaces.Module
	for _, model := range models {
		for _, module := range m.AvailableModules {
			if reflect.TypeOf(module).String() == model.Interface {
				modulesList = append(modulesList, module)
				break
			}
		}
	}

	return modulesList
}
func (m *moduleManagerImpl) Deactivated() []interfaces.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("flags & ? == 0", ACTIVE).Where("flags & ? != 0", INSTALLED).Get(&models)

	var modulesList []interfaces.Module
	for _, model := range models {
		for _, module := range m.AvailableModules {
			if reflect.TypeOf(module).String() == model.Interface {
				modulesList = append(modulesList, module)
				break
			}
		}
	}

	return modulesList
}
func (m *moduleManagerImpl) NotInstalled() []interfaces.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("(flags & ?) == 0", INSTALLED).Get(&models)

	var modulesList []interfaces.Module
	for _, module := range m.AvailableModules {
		found := false
		for _, model := range models {
			if reflect.TypeOf(module).String() == model.Interface {
				modulesList = append(modulesList, module)
				found = true
				break
			}
		}
		if !found {
			modulesList = append(modulesList, module)
		}
	}

	return modulesList
}

func (m *moduleManagerImpl) Install(module interfaces.Module) error {
	var model *ModuleModel
	err := m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if err != nil { panic(err) }

	err = m.services.Migrator().Migrate(module.Migrations())
	if err != nil { panic(err) }

	if model != nil {
		if model.Flags & INSTALLED != 0 {
			return ModuleError(module, "Target module is already installed")
		}

		model.Flags |= INSTALLED
		err = m.services.Repository().UpdateEntity(model)
		if err != nil { panic(err) }
	} else {
		model = &ModuleModel{
			Flags: INSTALLED,
			Interface: reflect.TypeOf(module).String(),
		}
		err := m.services.Repository().Insert(model)
		if err != nil { panic(err) }
	}

	if !module.Installed() {
		m.Purge(module)
		return ModuleError(module, "Module failed to install")
	}

	return nil
}
func (m *moduleManagerImpl) Activate(module interfaces.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	if model == nil || model.Flags & INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be activated")
	}
	if model.Flags & ACTIVE != 0 {
		return ModuleError(module, "Target module is already active")
	}

	m.services.Router().RegisterMiddleWares(module.MiddleWares())
	m.services.Router().RegisterGroups(module.GroupsHandlers())
	m.services.Router().RegisterRoutes(module.Routes())

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
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	m.services.Migrator().Rollback(module.Migrations())

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
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

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


