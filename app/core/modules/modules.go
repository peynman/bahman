package modules

import (
	"avalanche/app/core/app"
	"avalanche/app/core/interfaces"
	"avalanche/app/core/database"
)

func Initialize() {
	modules := app.InitAvalanchePlugins(app.ModulesPath("app"))

	for _, moduleInterface := range modules {
		module := moduleInterface.Interface().(interfaces.Module)

		database.DeployModule(module)
	}
}
