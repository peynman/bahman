package kernel

import "github.com/peyman-abdi/avalanche/app/interfaces/services"

type servicesImpl struct {
	app       services.Application
	config    services.Config
	logger    services.Logger
	repo      services.Repository
	mig       services.Migrator
	trans     services.Localization
	mm        services.ModuleManager
	router    services.Router
	templates services.RenderEngine
	redis     services.RedisClient
	cache     services.Cache
	names     map[string]interface{}
}

var instance *servicesImpl

func initializeServices(
	app services.Application,
	config services.Config,
	logger services.Logger,
	repo services.Repository,
	mig services.Migrator,
	trans services.Localization,
	mm services.ModuleManager,
	router services.Router,
	templates services.RenderEngine,
	redis services.RedisClient,
	cache services.Cache,
) services.Services {
	instance = new(servicesImpl)
	instance.names = make(map[string]interface{})

	instance.logger = logger
	instance.names["logger"] = logger
	instance.trans = trans
	instance.names["trans"] = trans
	instance.names["translations"] = trans
	instance.config = config
	instance.names["config"] = config
	instance.names["conf"] = config
	instance.app = app
	instance.names["app"] = app
	instance.names["application"] = app
	instance.mig = mig
	instance.names["migrations"] = mig
	instance.names["mig"] = mig
	instance.repo = repo
	instance.names["repositories"] = repo
	instance.names["repo"] = repo
	instance.mm = mm
	instance.names["modules"] = mm
	instance.router = router
	instance.names["router"] = router
	instance.templates = templates
	instance.names["templates"] = templates
	instance.redis = redis
	instance.names["redis"] = redis
	instance.cache = cache
	instance.names["cache"] = cache

	return instance
}

func (s *servicesImpl) Repository() services.Repository       { return s.repo }
func (s *servicesImpl) Migrator() services.Migrator           { return s.mig }
func (s *servicesImpl) Localization() services.Localization   { return s.trans }
func (s *servicesImpl) Config() services.Config               { return s.config }
func (s *servicesImpl) Logger() services.Logger               { return s.logger }
func (s *servicesImpl) Modules() services.ModuleManager       { return s.mm }
func (s *servicesImpl) App() services.Application             { return s.app }
func (s *servicesImpl) Router() services.Router               { return s.router }
func (s *servicesImpl) Renderer() services.RenderEngine { return s.templates }
func (s *servicesImpl) Redis() services.RedisClient           { return s.redis }
func (s *servicesImpl) GetByName(name string) interface{}     { return s.names[name] }
func (s *servicesImpl) Cache() services.Cache                 { return s.cache }
