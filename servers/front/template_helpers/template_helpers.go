package template_helpers

import (
	"os"
	"path/filepath"
)

func GetTemplates(dir string) []string {
	templates := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".gotmpl" {
			templates = append(templates, path)
		}
		return nil
	})
	return templates
}
