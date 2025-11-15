package github

import (
	"embed"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS
