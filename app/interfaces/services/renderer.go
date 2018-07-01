package services

import "io"

type Template struct {
	Name   string
	Path   string
}

type RenderEngine interface {
	AddTemplate(name string, content string) error
	ParseTemplates(templates []*Template) error
	ParseTemplate(template *Template) error
	Render(name string, params map[string]interface{}, writer io.Writer) error
}
