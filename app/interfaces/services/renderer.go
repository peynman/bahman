package services

import "io"

type Template struct {
	Name   string
	Path   string
}

type Parser interface {
	Load(services Services) error
	CanParse(template *Template) bool
	ParseFile(template *Template) error
	Parse(name string, content string) error
	Render(template *Template, params map[string]interface{}, writer io.Writer) error
}

type RenderEngine interface {
	Load(services Services) error
	AddParser(parser Parser) error
	AddTemplate(name string, content string, parser Parser) error
	ParseTemplates(templates []*Template) error
	ParseTemplate(template *Template) error
	Render(name string, params map[string]interface{}, writer io.Writer) error
}
