package core

import "io"

type Template struct {
	Name string
	Path string
	Weight int
}

type TemplateEngine interface {
	ParseTemplates(templates []*Template) error
	ParseTemplate(template *Template) error
	Render(name string, params map[string]interface{}, writer io.Writer) error
}

