package app_test

import (
	"github.com/peyman-abdi/avalanche/app/modules/services/app"
	"os"
	"path/filepath"
	"testing"
)

func TestPaths(t *testing.T) {
	application := app.Initialize(0, "test")

	root, err := os.Executable()
	if err != nil {
		panic(err)
	}
	root = filepath.Dir(root)

	testPath(filepath.Join(root, ""), application.RootPath(""), t)
	testPath(filepath.Join(root, "resources"), application.ResourcesPath(""), t)
	testPath(filepath.Join(root, "storage"), application.StoragePath(""), t)
	testPath(filepath.Join(root, "config"), application.ConfigPath(""), t)
	testPath(filepath.Join(root, "bin", "platforms", application.Platform(), application.Variant(), "modules"), application.ModulesPath(""), t)
}

func testPath(path1 string, path2 string, t *testing.T) {
	if path1 != path2 {
		t.Error("Folder was not correct,\nExpecting: " + path1 + "\nGot: " + path2)
	}
}
