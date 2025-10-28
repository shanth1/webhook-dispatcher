package templates

import (
	"html/template"
	"io/fs"

	"github.com/shanth1/gitrelay/templates"
)

func LoadTemplates() (*template.Template, error) {
	rootTmpl := template.New("root")

	err := fs.WalkDir(templates.TemplateFiles, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !fs.ValidPath(filePath) {
			return nil
		}

		content, err := fs.ReadFile(templates.TemplateFiles, filePath)
		if err != nil {
			return err
		}

		_, err = rootTmpl.New(filePath).Parse(string(content))
		return err
	})

	if err != nil {
		return nil, err
	}

	return rootTmpl, nil
}
