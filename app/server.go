package main

import (
	"github.com/sirupsen/logrus"
	"github.com/peyman-abdi/avalanche/app/core"
)

func main() {
	services := core.SetupKernel()
	application := services.App()

	services.Logger().InfoFields("Avalanche Server", logrus.Fields{
		"Version":   application.Version(),
		"BuildCode": application.BuildCode(),
		"Platform":  application.Platform(),
		"Variant":   application.Variant(),
		"BuildTime": application.BuildTime(),
	})
}


