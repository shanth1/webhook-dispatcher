package kanboard

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

func parseTemplates() (*template.Template, error) {
	return template.New("kanboard").ParseFS(templateFiles, "templates/*.tmpl")
}
