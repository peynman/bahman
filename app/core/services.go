package core

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

type servicesImpl struct {
	app interfaces.Application
	config interfaces.Config
	logger interfaces.Logger
	repo interfaces.Repository
	mig interfaces.Migrator
	trans interfaces.Localization
	mm interfaces.ModuleManager
	router interfaces.Router
}

var instance *servicesImpl

func initializeServices(
	app interfaces.Application,
	config interfaces.Config,
	logger interfaces.Logger,
	repo interfaces.Repository,
	mig interfaces.Migrator,
	trans interfaces.Localization,
	mm interfaces.ModuleManager,
	router interfaces.Router,
) interfaces.Services {
	instance = new(servicesImpl)

	instance.logger = logger
	instance.trans = trans
	instance.config = config
	instance.app = app
	instance.mig = mig
	instance.repo = repo
	instance.mm = mm
	instance.router = router

	return instance
}

func (s *servicesImpl) Repository() interfaces.Repository { return s.repo }
func (s *servicesImpl) Migrator() interfaces.Migrator { return s.mig }
func (s *servicesImpl) Localization() interfaces.Localization { return s.trans }
func (s *servicesImpl) Config() interfaces.Config { return s.config }
func (s *servicesImpl) Logger() interfaces.Logger { return s.logger }
func (s *servicesImpl) Modules() interfaces.ModuleManager { return s.mm }
func (s *servicesImpl) App() interfaces.Application { return s.app }
func (s *servicesImpl) Router() interfaces.Router { return s.router }
