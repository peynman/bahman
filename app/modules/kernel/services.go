package kernel

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/config"
	"github.com/peyman-abdi/bahman/app/modules/services/orm"
	"github.com/peyman-abdi/bahman/app/modules/services/logger"
	"github.com/peyman-abdi/bahman/app/modules/services/modules"
	"github.com/peyman-abdi/bahman/app/modules/services/http"
	"github.com/peyman-abdi/bahman/app/modules/services/renderer"
	"github.com/peyman-abdi/bahman/app/modules/services/trans"
	"github.com/peyman-abdi/bahman/app/modules/services/redis"
	"github.com/peyman-abdi/bahman/app/modules/services/cache"
	hash2 "github.com/peyman-abdi/bahman/app/modules/services/hash"
)

type servicesImpl struct {
	app    services.Application
	config services.Config
	logger services.Logger
	repo   services.Repository
	mig    services.Migratory
	trans  services.Localization
	mm     services.ModuleManager
	router services.Router
	ren    services.RenderEngine
	redis  services.RedisClient
	cache  services.Cache
	hash   services.Hash
	names  map[string]interface{}
}
func (s *servicesImpl) Repository() services.Repository     { return s.repo }
func (s *servicesImpl) Migrator() services.Migratory        { return s.mig }
func (s *servicesImpl) Localization() services.Localization { return s.trans }
func (s *servicesImpl) Config() services.Config             { return s.config }
func (s *servicesImpl) Logger() services.Logger             { return s.logger }
func (s *servicesImpl) Modules() services.ModuleManager     { return s.mm }
func (s *servicesImpl) App() services.Application           { return s.app }
func (s *servicesImpl) Router() services.Router             { return s.router }
func (s *servicesImpl) Renderer() services.RenderEngine 	  { return s.ren }
func (s *servicesImpl) Redis() services.RedisClient           { return s.redis }
func (s *servicesImpl) Hash() services.Hash           		  { return s.hash }
func (s *servicesImpl) GetByName(name string) interface{}     { return s.names[name] }
func (s *servicesImpl) Cache() services.Cache                 { return s.cache }

func initializeServices(application services.Application) services.Services {
	appConfig := config.New(application)
	appLogger := logger.New(appConfig)
	repo, migrator := orm.New(appConfig, appLogger)
	localizations := trans.New(appConfig, application, appLogger)
	hash := hash2.New(appConfig)
	mm := modules.New(appConfig, migrator)
	rendererRef := renderer.New(application, appLogger)
	redisRef := redis.New(appConfig, appLogger)
	routerRef := http.New(application, appConfig, appLogger, redisRef, rendererRef)
	cacheRef := cache.New(application, appConfig, appLogger, redisRef)

	instance := new(servicesImpl)
	instance.names = make(map[string]interface{})

	instance.logger = appLogger
	instance.names["logger"] = appLogger
	instance.trans = localizations
	instance.names["trans"] = localizations
	instance.names["translations"] = localizations
	instance.config = appConfig
	instance.names["config"] = appConfig
	instance.names["conf"] = appConfig
	instance.app = application
	instance.names["app"] = application
	instance.names["application"] = application
	instance.mig = migrator
	instance.names["migrations"] = migrator
	instance.names["mig"] = migrator
	instance.repo = repo
	instance.names["repositories"] = repo
	instance.names["repo"] = repo
	instance.mm = mm
	instance.names["modules"] = mm
	instance.router = routerRef
	instance.names["http"] = routerRef
	instance.ren = rendererRef
	instance.names["ren"] = rendererRef
	instance.redis = redisRef
	instance.names["redis"] = redisRef
	instance.cache = cacheRef
	instance.names["cache"] = cacheRef
	instance.hash = hash
	instance.names["hash"] = hash

	return instance
}
