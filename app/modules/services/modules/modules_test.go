package modules_test

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/bahman/app/modules/services/orm"
	"reflect"
	"testing"
	"github.com/peyman-abdi/avest"
)

var s services.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
}

func TestModuleStatus(t *testing.T) {
	testModule := new(avest.TestMigrationModule)

	mm := s.Modules()

	reflect.ValueOf(mm).Elem().FieldByName("AvailableModules").Set(
		reflect.ValueOf([]services.Module{testModule}),
	)
	repo := s.Repository()

	mm.Install(testModule)
	mm.Activate(testModule)

	var err error
	var migrations []*orm.MigrationModel
	if err = repo.Query(&orm.MigrationModel{}).GetAll(&migrations); err != nil {
		t.Error(err)
	}
	if len(migrations) != 1 {
		t.Errorf("Migrations not happend")
	}

	if len(mm.NotInstalled()) != 0 {
		t.Errorf("Module not installed list error")
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

	if err = repo.Query(&orm.MigrationModel{}).GetAll(&migrations); err != nil {
		t.Error(err)
	}
	if len(migrations) != 0 {
		t.Errorf("Migrations not rolledbacked")
	}

	if err = repo.Insert(&avest.TestMigrationModel{MyTest: "text"}); err != nil {
		t.Error(err)
	}
}
