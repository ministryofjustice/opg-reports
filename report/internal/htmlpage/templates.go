package htmlpage

import (
	"fmt"
	"path/filepath"
)

// GetTemplateFiles helper function to find all `.html` files
// within a given directory structure
func GetTemplateFiles(directory string) (files []string) {
	pattern := filepath.Join(directory, "**/**.html")
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("err:" + err.Error())
	}
	return
}
