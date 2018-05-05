package app

import (
	"os"
	"path/filepath"
)

var appRoot string

var (
	Version string
	Code string
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