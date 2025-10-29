package utils

import (
	"fmt"
	"path"
)

func GetTemplatePath(adapter, name string) string {
	return path.Join(adapter, fmt.Sprintf("%s.tmpl", name))
}
