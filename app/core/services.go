package core

import "github.com/peyman-abdi/avalanche/app/interfaces"

type servicesImpl struct {
	app interfaces.Application
	config interfaces.Config
	logger interfaces.Logger
	repo interfaces.Repository
	mig interfaces.Migrator
	trans interfaces.Localization
	mm interfaces.ModuleManager
}

func Initialize(
	app interfaces.Application,
	config interfaces.Config,
	logger interfaces.Logger,
	repo interfaces.Repository,
	mig interfaces.Migrator,
	trans interfaces.Localization,
	mm interfaces.ModuleManager,
) interfaces.Services {
	services := new(servicesImpl)
	services.logger = logger
	services.trans = trans
	services.config = config
	services.app = app
	services.mig = mig
	services.repo = repo
	services.mm = mm

	return services
}

func (s *servicesImpl) Repository() interfaces.Repository { return s.repo }
func (s *servicesImpl) Migrator() interfaces.Migrator { return s.mig }
func (s *servicesImpl) Localization() interfaces.Localization { return s.trans }
func (s *servicesImpl) Config() interfaces.Config { return s.config }
func (s *servicesImpl) Logger() interfaces.Logger { return s.logger }
func (s *servicesImpl) Modules() interfaces.ModuleManager { return s.mm }
func (s *servicesImpl) App() interfaces.Application { return s.app }
