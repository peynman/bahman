package main

import (
	"avalanche/app/core/logger"
	"avalanche/app/core/database"
	"github.com/sirupsen/logrus"
	"avalanche/app/core/config"
	"avalanche/app/core/app"
	"avalanche/app/core/trans"
	"avalanche/app/core/modules"
)

func main() {

	config.Initialize()
	logger.Initialize()
	database.Initialize()
	defer database.Close()
	trans.Initialize()

	modules.Initialize()

	logger.InfoFields("Avalanche Server", logrus.Fields{
		"version":  app.Version,
		"code":     app.Code,
		"platform": app.Platform,
		"variant":  app.Variant,
		"time":     app.BuildTime,
	})
}


