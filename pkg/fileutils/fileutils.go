// Package fileutils include extra features for handling files
//
// Includes helpers like copying a buffer to a file
package fileutils

import (
	"io"
	"os"
	"path/filepath"
)

// Copy copies the content of the source io.Reader the new destionation file path
func Copy(source io.Reader, destinationPath string) (destination *os.File, err error) {
	var directory string = filepath.Dir(destinationPath)

	// try to make the directory path
	if err = os.MkdirAll(directory, os.ModePerm); err != nil {
		return
	}
	// try to create the destination file
	if destination, err = os.Create(destinationPath); err != nil {
		return
	}
	defer destination.Close()
	// copy the source content to the destination
	_, err = io.Copy(destination, source)

	return
}
