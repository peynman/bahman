package template

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"io"
)

type templateHolder struct {
	meta     *core.Template
	template *jet.Template
}

type templateEngine struct {
	Views     *jet.Set
	logger    core.Logger
	templates []*templateHolder
}

func Initialize(app core.Application, logger core.Logger) core.TemplateEngine {
	t := new(templateEngine)
	t.logger = logger

	t.Views = jet.NewHTMLSet(app.ResourcesPath("views/templates"))

	if app.IsDebugMode() {
		t.Views.SetDevelopmentMode(true)
	}

	return t
}

func (t *templateEngine) ParseTemplates(templates []*core.Template) error {
	var err error
	for _, template := range templates {
		if err = t.ParseTemplate(template); err != nil {
			return err
		}
	}

	return nil
}
func (t *templateEngine) AddTemplate(name string, content string) error  {
	parsed, err := t.Views.LoadTemplate(name, content)
	if err != nil {
		t.logger.ErrorFields("Failed adding template", map[string]interface{}{
			"name": name,
			"content": content,
		})
		return err
	}

	t.templates = append(t.templates, &templateHolder{
		meta:     &core.Template{ Name: name, Path: name },
		template: parsed,
	})

	return nil
}
func (t *templateEngine) ParseTemplate(template *core.Template) error {
	parsed, err := t.Views.GetTemplate(template.Path)
	if err != nil {
		t.logger.ErrorFields("Failed loading template", map[string]interface{}{
			"name": template.Name,
			"path": template.Path,
		})
		return err
	}

	t.templates = append(t.templates, &templateHolder{
		meta:     template,
		template: parsed,
	})

	return nil
}

func (t *templateEngine) Render(name string, params map[string]interface{}, writer io.Writer) error {
	for _, temp := range t.templates {
		if temp.meta.Name == name {
			vars := make(jet.VarMap)
			for name, val := range params {
				vars.Set(name, val)
			}
			err := temp.template.Execute(writer, vars, nil)
			if err != nil {
				t.logger.ErrorFields("Could not render template due to error", map[string]interface{}{
					"name":  name,
					"error": err,
				})
				return err
			}
			return nil
		}
	}

	t.logger.ErrorFields("Could not find template to render", map[string]interface{}{
		"name": name,
	})
	return errors.New(fmt.Sprintf("Template with name %s not found", name))
}
