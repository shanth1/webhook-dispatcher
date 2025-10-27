package utils

import (
	"fmt"
	"path/filepath"
)

func GetTemplatePath(adapter, name string) string {
	return filepath.Join(adapter, fmt.Sprintf("%s.tmpl", name))
}
