package services

type Services interface {
	Repository() Repository
	Migrator() Migrator
	Localization() Localization
	Config() Config
	Logger() Logger
	Modules() ModuleManager
	App() Application
	Router() Router
	Renderer() RenderEngine
	Redis() RedisClient
	Cache() Cache
	GetByName(name string) interface{}
}

var Instance Services = nil