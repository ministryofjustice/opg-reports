package tmpl

import (
	"opg-reports/report/packages/fmtx"
	"path/filepath"
)

// Files helper function to find all `.html` files
// within a given directory structure
func Files(directory string) (templateFiles []string) {
	var patterns = []string{
		filepath.Join(directory, "*.html"),
		filepath.Join(directory, "**/**.html"),
		filepath.Join(directory, "**/**/**.html"),
		filepath.Join(directory, "**/**/**/**.html"),
	}

	for _, pattern := range patterns {
		if files, e := filepath.Glob(pattern); e == nil {
			templateFiles = append(templateFiles, files...)
		}
	}

	fmtx.Dump(directory)
	fmtx.Dump(templateFiles)

	return
}
