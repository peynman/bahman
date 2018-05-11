package app

import (
	"os"
	"path/filepath"
	"avalanche/app/core/interfaces"
	"strings"
	"plugin"
)

var appRoot string

var (
	Version string
	Code string
	VersionCode = 1
	Platform string
	Variant string
	BuildTime string
)

func init() {
	root, err := os.Executable()
	if err != nil {
		panic(err)
	}

	appRoot = filepath.Dir(root) + "/../../../../"
}

func StoragePath(path string) string {
	return filepath.Join(appRoot, "storage", path)
}

func ConfigPath(path string) string {
	return filepath.Join(appRoot, "config", path)
}

func RootPath(path string) string {
	return filepath.Join(appRoot, path)
}

func ModulesPath(path string) string {
	return filepath.Join(appRoot, "bin/platforms/", Platform, Variant, "modules", path)
}

func ResourcesPath(path string) string {
	return filepath.Join(appRoot, "resources", path)
}

func InitAvalanchePlugins(path string) []interfaces.AvalanchePlugin {
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

		if pluginInstance.AvalancheVersionCode() > VersionCode {
			panic("Plugin at path: " + moduleFile + " cannot run with this version of avalanche")
		}

		if pluginInstance.Initialize() {
			modules = append(modules, pluginInstance)
		}
	}

	return modules
}