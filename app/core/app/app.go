package app

import (
	"os"
	"path/filepath"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"strings"
	"plugin"
)

var (
	Version string
	Code string
	VersionCode = 1
	Platform string
	Variant string
	BuildTime string
)

type appImpl struct {
	appRoot string
}

func Initialize(roots int) interfaces.Application {
	root, err := os.Executable()
	if err != nil {
		panic(err)
	}

	root = filepath.Dir(root)
	for i := roots; i > 0; i-- {
		root = filepath.Join(root, "..")
	}

	app := &appImpl{
		appRoot: root,
	}

	return app
}

func (a *appImpl) Version() string {
	return Version
}
func (a *appImpl) Build() string {
	return Code
}
func (a *appImpl) BuildCode() int {
	return VersionCode
}
func (a *appImpl) Platform() string {
	return Platform
}
func (a *appImpl) Variant() string {
	return Variant
}
func (a *appImpl) BuildTime() string {
	return BuildTime
}


func (a *appImpl) StoragePath(path string) string {
	return filepath.Join(a.appRoot, "storage", path)
}

func (a *appImpl) ConfigPath(path string) string {
	return filepath.Join(a.appRoot, "config", path)
}

func (a *appImpl) RootPath(path string) string {
	return filepath.Join(a.appRoot, path)
}

func (a *appImpl) ModulesPath(path string) string {
	return filepath.Join(a.appRoot, "bin/platforms/", Platform, Variant, "modules", path)
}

func (a *appImpl) ResourcesPath(path string) string {
	return filepath.Join(a.appRoot, "resources", path)
}

func (a *appImpl) InitAvalanchePlugins(path string, services interfaces.Services) []interfaces.AvalanchePlugin {
	var moduleFiles []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		moduleFiles = append(moduleFiles, path)
		return nil
	})

	var modules []interfaces.AvalanchePlugin
	for _, moduleFile := range moduleFiles {
		if !strings.HasSuffix(moduleFile, ".so") {
			continue
		}

		pl, err := plugin.Open(moduleFile)
		if err != nil {
			panic(err)
		}

		pluginInstanceRef, err := pl.Lookup("PluginInstance")
		if err != nil {
			panic(err)
		}

		pluginInstance := *pluginInstanceRef.(*interfaces.AvalanchePlugin)
		if pluginInstance.Initialize(services) {
			modules = append(modules, pluginInstance)
		}
	}

	return modules
}