package renderer

import (
	"errors"
	"fmt"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"io"
	"github.com/peyman-abdi/bahman/app/modules/services/renderer/parsers"
	"reflect"
)

type templateHolder struct {
	meta *services.Template
	parser services.Parser
}

type renderEngine struct {
	logger    services.Logger
	parsers []services.Parser
	mainParser services.Parser
	templates map[string]*templateHolder
	instance services.Services
}

func New(app services.Application, logger services.Logger) services.RenderEngine {
	t := new(renderEngine)
	t.logger = logger
	t.parsers = []services.Parser{}
	t.mainParser = parsers.NewJetParser(app, logger)
	t.AddParser(t.mainParser)
	t.templates = make(map[string]*templateHolder)

	return t
}
func (t *renderEngine) Load(instance services.Services) error {
	t.instance = instance
	for _, parser := range t.parsers {
		if err := parser.Load(instance); err != nil {
			t.logger.ErrorFields("Failed loading parsers", map[string]interface{} {
				"parser": reflect.TypeOf(parser),
			})
			return err
		}
	}
	return nil
}

func (t *renderEngine) AddParser(parser services.Parser) error {
	t.parsers = append(t.parsers, parser)
	if t.instance != nil {
		if err := parser.Load(t.instance); err != nil {
			t.logger.ErrorFields("Failed loading parsers", map[string]interface{} {
				"parser": reflect.TypeOf(parser),
			})
			return err
		}
	}

	return nil
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
func (t *renderEngine) AddTemplate(name string, content string, parser services.Parser) error  {
	if parser == nil {
		parser = t.mainParser
	}
	err := parser.Parse(name, content)
	if err != nil {
		t.logger.ErrorFields("Failed parsing template", map[string]interface{} {
			"name": name,
			"content": content,
			"parser": reflect.TypeOf(parser),
		})
		return err
	}

	t.templates[name] = &templateHolder{&services.Template{Name:name, Path:""}, parser}
	return nil
}
func (t *renderEngine) ParseTemplate(template *services.Template) error {
	var err error
	for _, parser := range t.parsers {
		if parser.CanParse(template) {
			t.templates[template.Name] = &templateHolder{template, parser}
			if err := parser.ParseFile(template); err != nil {
				t.logger.ErrorFields("Failed parsing template", map[string]interface{} {
					"err": err,
					"name": template.Name,
					"path": template.Path,
					"parser": reflect.TypeOf(parser),
				})
				return err
			}
			return nil
		}
	}

	err = errors.New(fmt.Sprintf("Could not find parser to parse %s", template.Name))
	t.logger.ErrorFields("Could not find parser to parse", map[string]interface{}{
		"name": template.Name,
		"path": template.Path,
	})
	return err
}

func (t *renderEngine) Render(name string, params map[string]interface{}, writer io.Writer) error {
	for _,temp := range t.templates {
		if temp.meta.Name == name {
			for _, parser := range t.parsers {
				if parser == temp.parser {
					if err := parser.Render(temp.meta, params, writer); err != nil {
						t.logger.ErrorFields("Failed rendering template", map[string]interface{} {
							"error": err,
							"name": temp.meta.Name,
							"path": temp.meta.Path,
							"parser": reflect.TypeOf(parser),
						})
						return err
					}
					return nil
				}
			}
			return t.mainParser.Render(temp.meta, params, writer)
		}
	}
	t.logger.ErrorFields("Could not find renderer to render", map[string]interface{}{
		"name": name,
	})
	return errors.New(fmt.Sprintf("Could not find renderer to render template %s", name))
}
