package kernel

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/app"
)

func SetupServerServices() services.Services {
	return setupKernel(4, "SERVER")
}
func SetupCLIServices() services.Services {
	return setupKernel(4, "CLI")
}

func setupKernel(roots int, mode string) services.Services {
	application := app.New(roots, mode)
	s := initializeServices(application)
	services.Instance = s

	application.Load(s)
	s.Renderer().Load(s)
	s.Router().Load(s)

	s.Logger().LoadChannels(s)
	s.Modules().LoadModules(s)

	return s
}
