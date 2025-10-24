package templates

import (
	"html/template"

	"github.com/shanth1/telehook/templates"
)

func LoadTemplates() (*template.Template, error) {
	return template.ParseFS(templates.TemplateFiles, "*.tmpl")
}
