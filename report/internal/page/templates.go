package page

import (
	"path/filepath"
)

// GetTemplateFiles helper function to find all `.html` files
// within a given directory structure
func GetTemplateFiles(directory string) (templateFiles []string) {
	var patterns = []string{
		filepath.Join(directory, "*.html"),
		filepath.Join(directory, "**/**.html"),
	}

	for _, pattern := range patterns {
		if files, e := filepath.Glob(pattern); e == nil {
			templateFiles = append(templateFiles, files...)
		}
	}

	return
}
