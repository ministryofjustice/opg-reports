package page

import (
	"path/filepath"
)

// GetTemplateFiles helper function to find all `.html` files
// within a given directory structure
func GetTemplateFiles(directory string) (templateFiles []string) {
	var (
		pattern string
		files   = []string{}
	)
	// top level of dir
	pattern = filepath.Join(directory, "*.html")
	files, _ = filepath.Glob(pattern)
	templateFiles = append(templateFiles, files...)
	// sub levels
	pattern = filepath.Join(directory, "**/**.html")
	files, _ = filepath.Glob(pattern)
	templateFiles = append(templateFiles, files...)

	return
}
