package templates

import "embed"

//go:embed *.tmpl
var TemplateFiles embed.FS
