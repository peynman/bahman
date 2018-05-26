package interfaces

type Services interface {
	Repository() Repository
	Migrator() Migrator
	Localization() Localization
	Config() Config
	Logger() Logger
	Modules() ModuleManager
	App() Application
}
