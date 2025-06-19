package utils

import (
	"io"
	"os"
	"path/filepath"
)

// FileExists checks if the file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// FileCopy copies the content of the source io.Reader the new destionation file path
func FileCopy(source io.Reader, destinationPath string) (err error) {
	var directory string = filepath.Dir(destinationPath)
	var destination *os.File
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
