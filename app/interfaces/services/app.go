package services

type ApplicationError interface {
	error
	StackTrace() string
	File() string
	Line() int
	Type() string
	Code() int

	IsErrorType(t interface{}) bool
}

type Application interface {
	ConfigPath(path string) string
	StoragePath(path string) string
	AssetsPath(path string) string
	AssetsUrl(path string) string
	RootPath(path string) string
	ModulesPath(path string) string
	ResourcesPath(path string) string
	TemplatesPath(path string) string
	InitbahmanPlugins(path string, services Services) []Plugin
	Load(instance Services) error

	Error(t interface{}, code int, m string) ApplicationError

	Version() string
	Build() string
	BuildCode() int
	Platform() string
	Variant() string
	BuildTime() string
	Mode() string
	IsDebugMode() bool
}
