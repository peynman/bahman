package kernel

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/peyman-abdi/avalanche/app/modules/core/app"
	"github.com/peyman-abdi/avalanche/app/modules/core/config"
	"github.com/peyman-abdi/avalanche/app/modules/core/database"
	"github.com/peyman-abdi/avalanche/app/modules/core/logger"
	"github.com/peyman-abdi/avalanche/app/modules/core/modules"
	"github.com/peyman-abdi/avalanche/app/modules/core/router"
	"github.com/peyman-abdi/avalanche/app/modules/core/template"
	"github.com/peyman-abdi/avalanche/app/modules/core/trans"
)

func SetupServerKernel() core.Services {
	return setupKernel(4, "SERVER")
}
func SetupCLIKernel() core.Services {
	return setupKernel(4, "CLI")
}

func setupKernel(roots int, mode string) core.Services {
	application := app.Initialize(roots, mode)
	appConfig := config.Initialize(application)
	appLogger := logger.Initialize(appConfig)
	repo, migrator := database.Initialize(appConfig, appLogger)
	localizations := trans.Initialize(appConfig, application, appLogger)
	mm := modules.Initialize(appConfig, migrator)
	t := template.Initialize(application, appLogger)
	r := router.Initialize(appConfig, appLogger, t)

	s := initializeServices(
		application,
		appConfig,
		appLogger,
		repo,
		migrator,
		localizations,
		mm,
		r,
		t,
	)

	appLogger.LoadChannels(s)
	mm.LoadModules(s)

	return s
}
