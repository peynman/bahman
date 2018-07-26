package modules

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"reflect"
	"errors"
)

type moduleManagerImpl struct {
	AvailableModules []services.Module
	ModuleTableName  string
	services         services.Services
}

func New(config services.Config, migrator services.Migratory) services.ModuleManager {
	module := new(moduleManagerImpl)

	module.ModuleTableName = config.GetString("modules.table", "modules")

	if !migrator.HasTable(&ModuleModel{}) {
		migrator.AutoMigrate(&ModuleModel{})
	}

	return module
}

func (m *moduleManagerImpl) LoadModules(references services.Services) {
	modules := references.App().InitbahmanPlugins(references.App().ModulesPath("app"), references)

	m.services = references
	m.AvailableModules = make([]services.Module, len(modules))
	for index, moduleInterface := range modules {
		m.AvailableModules[index] = moduleInterface.Interface().(services.Module)
	}

	actives := m.Activated()
	for _, module := range actives {
		m.services.Router().RegisterMiddleWares(module.MiddleWares())
		m.services.Router().RegisterGroups(module.GroupsHandlers())
		m.services.Router().RegisterRoutes(module.Routes())

		if !module.Activated() {
			m.services.Router().RemoveMiddleWares(module.MiddleWares())
			m.services.Router().RemoveGroups(module.GroupsHandlers())
			m.services.Router().RemoveRoutes(module.Routes())
		}
	}
}

func (m *moduleManagerImpl) IsActive(module services.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if model != nil {
		return model.Flags&INSTALLED != 0 && model.Flags&ACTIVE != 0
	}

	return false
}
func (m *moduleManagerImpl) IsInstalled(module services.Module) bool {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if model != nil {
		return model.Flags&INSTALLED != 0
	}

	return false
}
func (m *moduleManagerImpl) List() []services.Module {
	return m.AvailableModules
}
func (m *moduleManagerImpl) Activated() []services.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("(flags & ?) != 0", ACTIVE).GetAll(&models)

	var modulesList []services.Module
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
func (m *moduleManagerImpl) Installed() []services.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("(flags & ?) != 0", INSTALLED).GetAll(&models)

	var modulesList []services.Module
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
func (m *moduleManagerImpl) Deactivated() []services.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("flags & ? == 0", ACTIVE).Where("flags & ? != 0", INSTALLED).GetAll(&models)

	var modulesList []services.Module
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
func (m *moduleManagerImpl) NotInstalled() []services.Module {
	var models []*ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("flags & ? != 0", INSTALLED).GetAll(&models)

	var modulesList []services.Module
	for _, module := range m.AvailableModules {
		installed := false
		for _, model := range models {
			if reflect.TypeOf(module).String() == model.Interface {
				installed = true
				break
			}
		}
		if !installed {
			modulesList = append(modulesList, module)
		}
	}

	return modulesList
}

func (m *moduleManagerImpl) Install(module services.Module) error {
	var model *ModuleModel
	err := m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
	if err != nil {
		panic(err)
	}

	err = m.services.Migrator().Migrate(module.Migrations())
	if err != nil {
		panic(err)
	}

	if model != nil {
		if model.Flags & INSTALLED != 0 {
			return ModuleError(module, "Target module is already installed")
		}

		model.Flags |= INSTALLED
		err = m.services.Repository().UpdateEntity(model)
		if err != nil {
			panic(err)
		}
	} else {
		model = &ModuleModel{
			Flags:     INSTALLED,
			Interface: reflect.TypeOf(module).String(),
		}
		err := m.services.Repository().Insert(model)
		if err != nil {
			panic(err)
		}
	}

	if !module.Installed() {
		m.Purge(module)
		return ModuleError(module, "Module failed to install")
	}

	return nil
}
func (m *moduleManagerImpl) SafeActivate(module services.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	if model == nil || model.Flags&INSTALLED == 0 {
		err := m.Install(module)
		if err != nil {
			return err
		}
		if model == nil {
			m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)
			if model == nil {
				return errors.New("")
			}
		}
	}
	if model.Flags&ACTIVE == 0 {
		if m.safeActivateModule(module) {
			model.Flags |= ACTIVE
			m.services.Repository().UpdateEntity(model)
		}
	}
	return nil
}
func (m *moduleManagerImpl) Activate(module services.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	if model == nil || model.Flags&INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be activated")
	}
	if model.Flags&ACTIVE != 0 {
		return ModuleError(module, "Target module is already active")
	}

	if m.safeActivateModule(module) {
		model.Flags |= ACTIVE
		m.services.Repository().UpdateEntity(model)
	}

	return nil
}
func (m *moduleManagerImpl) Purge(module services.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	if model == nil || model.Flags & INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be deactivated")
	}

	if model.Flags & ACTIVE != 0 {
		m.safeDeactivateModule(module)
	}
	m.services.Migrator().Rollback(module.Migrations())

	model.Flags = 0
	m.services.Repository().UpdateEntity(model)

	module.Purged()

	return nil
}
func (m *moduleManagerImpl) Deactivate(module services.Module) error {
	var model *ModuleModel
	m.services.Repository().Query(&ModuleModel{}).Where("interface = ?", reflect.TypeOf(module).String()).GetFirst(&model)

	if model == nil || model.Flags&INSTALLED == 0 {
		return ModuleError(module, "Target module is not installed yet and con not be deactivated")
	}
	if model.Flags&ACTIVE == 0 {
		return ModuleError(module, "Target module is not activated yet and can not be deactivated")
	}

	m.safeDeactivateModule(module)
	model.Flags &= ^ACTIVE
	m.services.Repository().UpdateEntity(model)

	return nil
}
