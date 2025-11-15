// internal/templates/registry.go
package templates

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"
)

type Source struct {
	FS       fs.FS
	Patterns []string
}

type Registry struct {
	sources map[string]Source
}

func NewRegistry() *Registry {
	return &Registry{
		sources: make(map[string]Source),
	}
}

func (r *Registry) RegisterSource(name string, source Source) error {
	if _, exists := r.sources[name]; exists {
		return fmt.Errorf("template source '%s' already registered", name)
	}
	r.sources[name] = source
	return nil
}

func (r *Registry) LoadAll(funcMap template.FuncMap) (*template.Template, error) {
	rootTmpl := template.New("root").Funcs(funcMap)

	for sourceName, source := range r.sources { // sourceName - это "github", "kanboard" и т.д.
		for _, pattern := range source.Patterns {
			matches, err := fs.Glob(source.FS, pattern)
			if err != nil {
				return nil, fmt.Errorf("failed to glob pattern '%s' in source '%s': %w", pattern, sourceName, err)
			}

			if len(matches) == 0 {
				continue
			}

			for _, matchPath := range matches {

				content, err := fs.ReadFile(source.FS, matchPath)
				if err != nil {
					return nil, fmt.Errorf("failed to read template file '%s' from source '%s': %w", matchPath, sourceName, err)
				}

				fileName := filepath.Base(matchPath)

				templateFullName := fmt.Sprintf("%s/%s", sourceName, fileName)

				_, err = rootTmpl.New(templateFullName).Parse(string(content))
				if err != nil {
					return nil, fmt.Errorf("failed to parse template '%s': %w", templateFullName, err)
				}
			}
		}
	}
	return rootTmpl, nil
}
