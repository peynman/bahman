package modules_test

import (
	"github.com/peyman-abdi/avalanche/app/modules/core/database"
	"github.com/peyman-abdi/testil"
	"testing"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"reflect"
)

var envs = map[string]string{}
var configs = map[string]interface{}{
	"app.hjson": map[string]interface{}{
		"debug": true,
	},
	"database.hjson": map[string]interface{}{
		"app": "sqlite3",
		"runtime": map[string]interface{}{
			"migrations": "migrations",
			"connection": "sqlite3",
		},
		"connections": map[string]interface{}{
			"sqlite3": map[string]interface{}{
				"driver": "sqlite3",
				"file":   "storage(\"test.db\")",
			},
		},
	},
}

var s core.Services
func init() {
	s = testil.MockServices(configs, envs)
}

func TestModuleStatus(t *testing.T) {
	testModule := new(testil.TestMigrationModule)

	mm := s.Modules()

	reflect.ValueOf(mm).Elem().FieldByName("AvailableModules").Set(
		reflect.ValueOf([]core.Module{testModule}),
	)
	repo := s.Repository()

	mm.Install(testModule)
	mm.Activate(testModule)

	var err error
	var migrations []*database.MigrationModel
	if err = repo.Query(&database.MigrationModel{}).Get(&migrations); err != nil {
		t.Error(err)
	}
	if len(migrations) != 1 {
		t.Errorf("Migrations not happend")
	}

	if !mm.IsInstalled(testModule) {
		t.Errorf("Module not installed!")
	}
	if !mm.IsActive(testModule) {
		t.Errorf("Module not active")
	}

	if c := len(mm.Activated()); c != 1 {
		t.Errorf("Activated list faled: %d", c)
	}

	mm.Deactivate(testModule)

	if mm.IsActive(testModule) {
		t.Errorf("Module not deactive")
	}

	if len(mm.Activated()) != 0 {
		t.Errorf("Activated module faled")
	}
	if len(mm.Deactivated()) != 1 {
		t.Errorf("Deactivated list faled")
	}
	if len(mm.Installed()) != 1 {
		t.Errorf("Installed list failed")
	}

	mm.Purge(testModule)
	if len(mm.NotInstalled()) != 1 {
		t.Errorf("Not installed list failed")
	}

	if err = repo.Query(&database.MigrationModel{}).Get(&migrations); err != nil {
		t.Error(err)
	}
	if len(migrations) != 0 {
		t.Errorf("Migrations not rolledbacked")
	}

	if err = repo.Insert(&testil.TestMigrationModel{MyTest: "text"}); err != nil {
		t.Error(err)
	}
}
