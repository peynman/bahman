package template_test

import (
	"bytes"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"strings"
	"testing"
	"github.com/peyman-abdi/avest"
)

var s core.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
	avest.CreateTemplateFiles(s.App(), avest.SimpleTemplates)
}

func TestInitialize(t *testing.T) {
	module := new(avest.TestRouteModule)
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

	avest.TestHTMLRequest(t, "http://127.0.0.1:8181/home", "This content will be yielded in the layout above.")
	avest.TestHTMLRequest(t, "http://127.0.0.1:8181/error", "internal server error")
}
