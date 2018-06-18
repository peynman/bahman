package main

import (
	"github.com/peyman-abdi/avalanche/app/modules/core"
	"github.com/sirupsen/logrus"
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

	services.Router().Serve()
}
