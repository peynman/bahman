package router_test

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"os"
	"github.com/peyman-abdi/testil"
	"github.com/peyman-abdi/avalanche/app/core/config"
	"github.com/peyman-abdi/avalanche/app/core/logger"
	"github.com/peyman-abdi/avalanche/app/core/database"
	"github.com/peyman-abdi/avalanche/app/core/router"
	application "github.com/peyman-abdi/avalanche/app/core/app"
	"github.com/peyman-abdi/avalanche/app/core/modules"
	"testing"
	"time"
)

var app interfaces.Application
var conf interfaces.Config
var log interfaces.Logger
var repo interfaces.Repository
var mig interfaces.Migrator
var mm interfaces.ModuleManager
var r interfaces.Router
var s = new(ServicesMock)

type ServicesMock struct {
}
func (s *ServicesMock) Repository() interfaces.Repository { return repo }
func (s *ServicesMock) Migrator() interfaces.Migrator { return mig }
func (s *ServicesMock) Localization() interfaces.Localization { return nil }
func (s *ServicesMock) Config() interfaces.Config { return conf }
func (s *ServicesMock) Logger() interfaces.Logger { return log }
func (s *ServicesMock) Modules() interfaces.ModuleManager { return mm }
func (s *ServicesMock) App() interfaces.Application { return app }
func (s *ServicesMock) Router() interfaces.Router { return r }

var envs = map[string]string {
}
var configs = map[string]interface{} {
	"app.hjson": map[string]interface{} {
		"debug": true,
	},
	"database.hjson": map[string]interface{} {
		"app": "sqlite3",
		"runtime": map[string]interface{} {
			"migrations": "migrations",
			"connection": "sqlite3",
		},
		"connections": map[string]interface{} {
			"sqlite3": map[string]interface{} {
				"driver": "sqlite3",
				"file": "storage(\"test.db\")",
			},
		},
	},
	"server.hjson": map[string]interface{} {
		"address": "127.0.0.1",
		"port": 8181,
	},
}

func init()  {
	app = application.Initialize(0)
	os.MkdirAll(app.StoragePath(""), 0700)

	testil.CreateConfigFiles(app, configs)

	conf = config.Initialize(app)
	log = logger.Initialize(conf)
	log.LoadConsole()

	repo, mig = database.Initialize(conf, log)

	r = router.Initialize(conf, log)

	mm = modules.Initialize(conf, mig)
	mm.LoadModules(s)
}
func TestRoutes(t *testing.T)  {
	module := new(testil.TestMigrationModule)
	module.S = s

	mig.Migrate(module.Migrations())
	r.RegisterGroups(module.GroupsHandlers())
	r.RegisterMiddleWares(module.MiddleWares())
	r.RegisterRoutes(module.Routes())

	ch := make(chan error)
	go func() {
		ch <- r.Serve()
	}()

	testil.TestGetRequest(t, "http://127.0.0.1:8181/", "hello world")
	testil.TestGetRequest(t, "http://127.0.0.1:8181/api", "hello api")
	testil.TestGetRequest(t, "http://127.0.0.1:8181/not-exist", "Page not found!")
	testil.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/tests/id/20/str/avalanche", map[string]interface{} {
		"route:test": 7,
		"middleware:auth": 6,
		"group:tests": 5,
		"group:api": 4,
		"id": "20",
		"name": "avalanche",
	})
	testil.TestPostJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{} {
		"text": "something",
	}, map[string]interface{} {
		"ID": 1,
		"MyTest": "something",
		"MyInt": nil,
	})

	testil.TestPostJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{} {
		"text": "something 2",
		"int": 35,
	}, map[string]interface{} {
		"ID": 2,
		"MyTest": "something 2",
		"MyInt": 35,
	})
	testil.TestPutJSONRequest(t, "http://127.0.0.1:8181/api/models/1", map[string]interface{} {
		"text": "new things",
		"int": 33,
	}, map[string]interface{} {
		"ID": 1,
		"MyTest": "new things",
		"MyInt": 33,
	})
	testil.TestPutJSONRequestString(t, "http://127.0.0.1:8181/api/models/3", map[string]interface{} {
	}, "object not found")
	testil.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{} {
		"count": "2",
	})
	testil.TestDeleteRequestString(t, "http://127.0.0.1:8181/api/models/1", "deleted")
	testil.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{} {
		"count": "1",
	})

	select {
	case err := <- ch:
		if err != nil {
			t.Error(err)
		}
		case <- time.After(100 * time.Millisecond):
			return
	}


}

