package main

import (
	"github.com/peyman-abdi/avalanche/app/modules/kernel"
	"github.com/sirupsen/logrus"
	"github.com/peyman-abdi/avalanche/app/modules/auth"
)

func main() {
	services := kernel.SetupServerKernel()
	application := services.App()

	services.Modules().SafeActivate(new (auth.AuthenticationModule))

	services.Logger().InfoFields("Avalanche Server", logrus.Fields{
		"Version":   application.Version(),
		"BuildCode": application.BuildCode(),
		"Platform":  application.Platform(),
		"Variant":   application.Variant(),
		"BuildTime": application.BuildTime(),
	})

	activesModules := services.Modules().Activated()
	for _, module := range activesModules {
		services.Logger().InfoFields("Using Module", map[string]interface{}{
			"module": module.Title(),
		})
	}

	services.Router().Serve()
}
