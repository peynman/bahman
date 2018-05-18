package main

import (
	"github.com/peyman-abdi/avalanche/app/core/logger"
	"github.com/peyman-abdi/avalanche/app/core/database"
	"github.com/sirupsen/logrus"
	"github.com/peyman-abdi/avalanche/app/core/config"
	"github.com/peyman-abdi/avalanche/app/core/app"
	"github.com/peyman-abdi/avalanche/app/core/trans"
	"github.com/peyman-abdi/avalanche/app/core/modules"
	"github.com/peyman-abdi/avalanche/app/core"
)

func main() {
	application := app.Initialize()
	appConfig := config.Initialize(application)
	appLogger := logger.Initialize()
	repo, migrator := database.Initialize(appConfig)
	defer database.Close()
	localizations := trans.Initialize(appConfig, application, appLogger)

	mm := modules.Initialize(appConfig)

	s := core.Initialize(application, appConfig, appLogger, repo, migrator, localizations, mm)

	appLogger.LoadChannels(s)
	mm.LoadModules(s)

	appLogger.InfoFields("Avalanche Server", logrus.Fields{
		"Version":   application.Version(),
		"BuildCode": application.BuildCode(),
		"Platform":  application.Platform(),
		"Variant":   application.Variant(),
		"BuildTime": application.BuildTime(),
	})
}


