package interfaces

type Module interface {
	Migrations() []Migratable

	Activated() bool
	Installed() bool
	Deactivated()
	Purged()
}

type ModuleManager interface {
	LoadModules(Services)
	List() []Module
	IsActive(Module) bool
	IsInstalled(Module) bool
	Install(Module) error
	Activate(Module) error
	Purge(Module) error
	Deactivate(Module) error
}