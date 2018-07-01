package renderer

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"io"
)

type templateHolder struct {
	meta     *services.Template
	template *jet.Template
}

type renderEngine struct {
	Views     *jet.Set
	logger    services.Logger
	templates []*templateHolder
}

func Initialize(app services.Application, logger services.Logger) services.RenderEngine {
	t := new(renderEngine)
	t.logger = logger

	t.Views = jet.NewHTMLSet(app.TemplatesPath(""))

	if app.IsDebugMode() {
		t.Views.SetDevelopmentMode(true)
	}

	return t
}

func (t *renderEngine) ParseTemplates(templates []*services.Template) error {
	var err error
	for _, template := range templates {
		if err = t.ParseTemplate(template); err != nil {
			return err
		}
	}

	return nil
}
func (t *renderEngine) AddTemplate(name string, content string) error  {
	parsed, err := t.Views.LoadTemplate(name, content)
	if err != nil {
		t.logger.ErrorFields("Failed adding renderer", map[string]interface{}{
			"name": name,
			"content": content,
		})
		return err
	}

	t.templates = append(t.templates, &templateHolder{
		meta:     &services.Template{ Name: name, Path: name },
		template: parsed,
	})

	return nil
}
func (t *renderEngine) ParseTemplate(template *services.Template) error {
	parsed, err := t.Views.GetTemplate(template.Path)
	if err != nil {
		t.logger.ErrorFields("Failed loading renderer", map[string]interface{}{
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

func (t *renderEngine) Render(name string, params map[string]interface{}, writer io.Writer) error {
	for _, temp := range t.templates {
		if temp.meta.Name == name {
			vars := make(jet.VarMap)
			for name, val := range params {
				vars.Set(name, val)
			}
			err := temp.template.Execute(writer, vars, nil)
			if err != nil {
				t.logger.ErrorFields("Could not render renderer due to error", map[string]interface{}{
					"name":  name,
					"error": err,
				})
				return err
			}
			return nil
		}
	}

	t.logger.ErrorFields("Could not find renderer to render", map[string]interface{}{
		"name": name,
	})
	return errors.New(fmt.Sprintf("Template with name %s not found", name))
}
