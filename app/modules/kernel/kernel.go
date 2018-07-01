package kernel

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/peyman-abdi/avalanche/app/modules/services/app"
	"github.com/peyman-abdi/avalanche/app/modules/services/config"
	"github.com/peyman-abdi/avalanche/app/modules/services/database"
	"github.com/peyman-abdi/avalanche/app/modules/services/logger"
	"github.com/peyman-abdi/avalanche/app/modules/services/modules"
	"github.com/peyman-abdi/avalanche/app/modules/services/router"
	"github.com/peyman-abdi/avalanche/app/modules/services/renderer"
	"github.com/peyman-abdi/avalanche/app/modules/services/trans"
	"github.com/peyman-abdi/avalanche/app/modules/services/redis"
	"github.com/peyman-abdi/avalanche/app/modules/services/cache"
)

func SetupServerKernel() services.Services {
	return setupKernel(4, "SERVER")
}
func SetupCLIKernel() services.Services {
	return setupKernel(4, "CLI")
}

func setupKernel(roots int, mode string) services.Services {
	application := app.Initialize(roots, mode)
	appConfig := config.Initialize(application)
	appLogger := logger.Initialize(appConfig)
	repo, migrator := database.Initialize(appConfig, appLogger)
	localizations := trans.Initialize(appConfig, application, appLogger)
	mm := modules.Initialize(appConfig, migrator)
	t := renderer.Initialize(application, appLogger)
	r := router.Initialize(application, appConfig, appLogger, t)
	rredis := redis.Initialize(appConfig)
	c := cache.Initialize(application, appConfig, appLogger, rredis)

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
		rredis,
		c,
	)
	services.Instance = s

	appLogger.LoadChannels(s)
	mm.LoadModules(s)
	return s
}
