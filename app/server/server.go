package main

import (
	"github.com/peyman-abdi/bahman/app/modules/kernel"
	"github.com/peyman-abdi/bahman/app/modules/auth"
)

func main() {
	services := kernel.SetupServerServices()

	services.Modules().SafeActivate(new (auth.AuthenticationModule))

	services.Logger().InfoFields("Server ready", map[string]interface{} {
		"Version":   services.App().Version(),
		"BuildCode": services.App().BuildCode(),
		"Platform":  services.App().Platform(),
		"Variant":   services.App().Variant(),
		"BuildTime": services.App().BuildTime(),
	})

	activesModules := services.Modules().Activated()
	for _, module := range activesModules {
		services.Logger().InfoFields("Modules", map[string]interface{}{
			"module": module.Title(),
		})
	}

	services.Router().Serve()
}
