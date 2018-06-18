package core

import (
	"github.com/peyman-abdi/avalanche/app/modules/core/app"
	"github.com/peyman-abdi/avalanche/app/modules/core/config"
	"github.com/peyman-abdi/avalanche/app/modules/core/database"
	"github.com/peyman-abdi/avalanche/app/modules/core/logger"
	"github.com/peyman-abdi/avalanche/app/modules/core/modules"
	"github.com/peyman-abdi/avalanche/app/modules/core/router"
	"github.com/peyman-abdi/avalanche/app/modules/core/trans"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

func SetupKernel() core.Services {
	application := app.Initialize(4)
	appConfig := config.Initialize(application)
	appLogger := logger.Initialize(appConfig)
	repo, migrator := database.Initialize(appConfig, appLogger)
	localizations := trans.Initialize(appConfig, application, appLogger)
	mm := modules.Initialize(appConfig, migrator)
	r := router.Initialize(appConfig, appLogger)

	s := initializeServices(
		application,
		appConfig,
		appLogger,
		repo,
		migrator,
		localizations,
		mm,
		r,
	)

	appLogger.LoadChannels(s)
	mm.LoadModules(s)

	return s
}
