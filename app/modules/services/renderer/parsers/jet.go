package parsers

import (
	"github.com/CloudyKit/jet"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"io"
	"strings"
	"errors"
	"fmt"
)

type jetTemplateHolder struct {
	meta     *services.Template
	template *jet.Template
}

type jetParser struct {
	Views     *jet.Set
	logger    services.Logger
	instance  services.Services
	templates map[string]*jetTemplateHolder
}

func (j *jetParser) CanParse(template *services.Template) bool {
	return strings.HasSuffix(template.Path, ".jet")
}

func (t *jetParser) ParseFile(template *services.Template) error {
	parsed, err := t.Views.GetTemplate(template.Path)
	if err != nil {
		t.logger.ErrorFields("Failed parsing jet template", map[string]interface{}{
			"error": err,
			"name": template.Name,
			"path": template.Path,
		})
		return err
	}

	t.templates[template.Name] = &jetTemplateHolder{
		meta:     template,
		template: parsed,
	}

	return nil
}

func (t *jetParser) Parse(name string, content string) error {
	parsed, err := t.Views.LoadTemplate(name, content)
	if err != nil {
		t.logger.ErrorFields("Failed parsing jet template content", map[string]interface{}{
			"error": err,
			"name": name,
			"content": content,
		})
		return err
	}

	t.templates[name] = &jetTemplateHolder{
		meta:     &services.Template{ Name: name, Path: name },
		template: parsed,
	}
	return nil
}

func (t *jetParser) Render(template *services.Template, params map[string]interface{}, writer io.Writer) error {
	temp := t.templates[template.Name]
	if temp != nil {
		vars := make(jet.VarMap)
		for name, val := range params {
			vars.Set(name, val)
		}
		err := temp.template.Execute(writer, vars, t.instance)
		if err != nil {
			t.logger.ErrorFields("Could not render renderer due to error", map[string]interface{}{
				"name":  template.Name,
				"error": err,
			})
			return err
		}
		return nil
	}

	return errors.New(fmt.Sprintf("Template with name %s not found", template.Name))
}

func (t *jetParser) Load(services services.Services) error {
	t.instance = services
	return nil
}

func NewJetParser(app services.Application, logger services.Logger) services.Parser {
	t := new(jetParser)
	t.Views = jet.NewHTMLSet(app.TemplatesPath(""))
	t.logger = logger
	t.templates = make(map[string]*jetTemplateHolder)

	if app.IsDebugMode() {
		t.Views.SetDevelopmentMode(true)
	}

	return t
}
