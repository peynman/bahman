package modules

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
)

func (m *moduleManagerImpl) safeActivateModule(module services.Module) bool {
	if module.MiddleWares() != nil {
		if err := m.services.Router().RegisterMiddleWares(module.MiddleWares()); err != nil {
			m.services.Logger().ErrorFields("Could not register middle wares of module", map[string]interface{} {
				"err": err,
				"module": module,
			})
		}
	}
	if module.GroupsHandlers() != nil {
		if err := m.services.Router().RegisterGroups(module.GroupsHandlers()); err != nil {
			m.services.Logger().ErrorFields("Could not register group handlers of module", map[string]interface{} {
				"err": err,
				"module": module,
			})
		}
	}
	if module.Routes() != nil {
		if err := m.services.Router().RegisterRoutes(module.Routes()); err != nil {
			m.services.Logger().ErrorFields("Could not register routes of module", map[string]interface{} {
				"err": err,
				"module": module,
			})
		}
	}
	if module.Templates() != nil {
		if err := m.services.Renderer().ParseTemplates(module.Templates()); err != nil {
			m.services.Logger().ErrorFields("Could not register templates of module", map[string]interface{} {
				"err": err,
				"module": module,
			})
		}
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
