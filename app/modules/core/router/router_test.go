package router_test

import (
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"testing"
	"time"
	"github.com/peyman-abdi/avest"
)

var s core.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
	avest.CreateTemplateFiles(s.App(), avest.SimpleTemplates)
}

func TestRoutes(t *testing.T) {
	module := new(avest.TestRouteModule)
	module.S = s

	s.Modules().Install(module)
	s.Modules().Activate(module)

	ch := make(chan error)
	go func() {
		ch <- s.Router().Serve()
	}()

	avest.TestGetRequest(t, "http://127.0.0.1:8181/", "hello world")
	avest.TestGetRequest(t, "http://127.0.0.1:8181/api", "hello api")
	avest.TestGetRequest(t, "http://127.0.0.1:8181/not-exist", "Page not found!")
	avest.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/tests/id/20/str/avalanche", map[string]interface{}{
		"route:test":      7,
		"middleware:auth": 6,
		"group:tests":     5,
		"group:api":       4,
		"id":              "20",
		"name":            "avalanche",
	})
	avest.TestPostJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{}{
		"text": "something",
	}, map[string]interface{}{
		"ID":     1,
		"MyTest": "something",
		"MyInt":  nil,
	})

	avest.TestPostJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{}{
		"text": "something 2",
		"int":  35,
	}, map[string]interface{}{
		"ID":     2,
		"MyTest": "something 2",
		"MyInt":  35,
	})
	avest.TestPutJSONRequest(t, "http://127.0.0.1:8181/api/models/1", map[string]interface{}{
		"text": "new things",
		"int":  33,
	}, map[string]interface{}{
		"ID":     1,
		"MyTest": "new things",
		"MyInt":  33,
	})
	avest.TestPutJSONRequestString(t, "http://127.0.0.1:8181/api/models/3", map[string]interface{}{}, "object not found")
	avest.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{}{
		"count": "2",
	})
	avest.TestDeleteRequestString(t, "http://127.0.0.1:8181/api/models/1", "deleted")
	avest.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/models", map[string]interface{}{
		"count": "1",
	})

	select {
	case err := <-ch:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(100 * time.Millisecond):
		return
	}

}
