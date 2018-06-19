package kernel

import "github.com/peyman-abdi/avalanche/app/interfaces/core"

type servicesImpl struct {
	app    core.Application
	config core.Config
	logger core.Logger
	repo   core.Repository
	mig    core.Migrator
	trans  core.Localization
	mm     core.ModuleManager
	router core.Router
	templates core.TemplateEngine
}

var instance *servicesImpl

func initializeServices(
	app core.Application,
	config core.Config,
	logger core.Logger,
	repo core.Repository,
	mig core.Migrator,
	trans core.Localization,
	mm core.ModuleManager,
	router core.Router,
	templates core.TemplateEngine,
) core.Services {
	instance = new(servicesImpl)

	instance.logger = logger
	instance.trans = trans
	instance.config = config
	instance.app = app
	instance.mig = mig
	instance.repo = repo
	instance.mm = mm
	instance.router = router
	instance.templates = templates

	return instance
}

func (s *servicesImpl) Repository() core.Repository     { return s.repo }
func (s *servicesImpl) Migrator() core.Migrator         { return s.mig }
func (s *servicesImpl) Localization() core.Localization { return s.trans }
func (s *servicesImpl) Config() core.Config             { return s.config }
func (s *servicesImpl) Logger() core.Logger             { return s.logger }
func (s *servicesImpl) Modules() core.ModuleManager     { return s.mm }
func (s *servicesImpl) App() core.Application           { return s.app }
func (s *servicesImpl) Router() core.Router             { return s.router }
func (s *servicesImpl) TemplateEngine() core.TemplateEngine   { return s.templates }
