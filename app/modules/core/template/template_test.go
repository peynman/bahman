package template_test

import (
	"bytes"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/peyman-abdi/testil"
	"strings"
	"testing"
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
	"server.hjson": map[string]interface{}{
		"address": "127.0.0.1",
		"port":    8181,
	},
}

var s core.Services

func init() {
	s = testil.MockServices(configs, envs)
	testil.CreateTemplateFiles(s.App(), testil.SimpleTemplates)
}

func TestInitialize(t *testing.T) {
	module := new(testil.TestRouteModule)
	module.S = s

	s.Modules().Install(module)
	s.Modules().Activate(module)

	buf := bytes.NewBufferString("")
	err := s.TemplateEngine().Render("home", nil, buf)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(buf.String(), "This content will be yielded in the layout above.") {
		t.Errorf("template not rendered as expected: %s", buf.String())
	}

	err = s.TemplateEngine().Render("not found", nil, buf)
	if err == nil {
		t.Errorf("template should not be found but no errors returned")
	}

	ch := make(chan error)
	go func() {
		ch <- s.Router().Serve()
	}()

	testil.TestHTMLRequest(t, "http://127.0.0.1:8181/home", "This content will be yielded in the layout above.")
	testil.TestHTMLRequest(t, "http://127.0.0.1:8181/error", "internal server error")
}
