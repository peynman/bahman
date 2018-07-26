package http_test

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"testing"
	"github.com/peyman-abdi/avest"
	"time"
)

var s services.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
	avest.CreateTemplateFiles(s.App(), avest.SimpleTemplates)
}

func TestRoutes(t *testing.T) {
	module := new(avest.TestRouteModule)
	module.S = s

	err := s.Modules().Install(module)
	if err != nil {
		t.Error(err)
	}
	err = s.Modules().Activate(module)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan error)
	go func() {
		ch <- s.Router().Serve()
	}()

	avest.TestGetJSONRequest(t, "http://127.0.0.1:8181/api/tests/id/20/str/bahman", map[string]interface{}{
		"middleware:auth": 1,
		"group:api":       2,
		"group:tests":     3,
		"route:test":      4,
		"id":              "20",
		"name":            "bahman",
	})
	avest.TestGetRequest(t, "http://127.0.0.1:8181/", "hello world")
	avest.TestSession(t,
		"http://127.0.0.1:8181/api/tests/session/set",
		"http://127.0.0.1:8181/api/tests/session/get",
		map[string]string {
			"handle": "route",
			"value": "test",
		})
	avest.TestGetRequest(t, "http://127.0.0.1:8181/not-exist", "Page not found!")
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
