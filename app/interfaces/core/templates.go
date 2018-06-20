package core

import "io"

type Template struct {
	Name   string
	Path   string
}

type TemplateEngine interface {
	AddTemplate(name string, content string) error
	ParseTemplates(templates []*Template) error
	ParseTemplate(template *Template) error
	Render(name string, params map[string]interface{}, writer io.Writer) error
}
