package modules

import "github.com/peyman-abdi/avalanche/app/interfaces/services"

func (m *moduleManagerImpl) safeActivateModule(module services.Module) bool {
	if module.MiddleWares() != nil {
		m.services.Router().RegisterMiddleWares(module.MiddleWares())
	}
	if module.GroupsHandlers() != nil {
		m.services.Router().RegisterGroups(module.GroupsHandlers())
	}
	if module.Routes() != nil {
		m.services.Router().RegisterRoutes(module.Routes())
	}
	if module.Templates() != nil {
		m.services.Renderer().ParseTemplates(module.Templates())
	}

	if !module.Activated() {
		m.safeDeactivateModule(module)
		return false
	}

	return true
}

func (m *moduleManagerImpl) safeDeactivateModule(module services.Module) {
	if module.MiddleWares() != nil {
		m.services.Router().RemoveMiddleWares(module.MiddleWares())
	}
	if module.GroupsHandlers() != nil {
		m.services.Router().RemoveGroups(module.GroupsHandlers())
	}
	if module.Routes() != nil {
		m.services.Router().RemoveRoutes(module.Routes())
	}

	module.Deactivated()
}
