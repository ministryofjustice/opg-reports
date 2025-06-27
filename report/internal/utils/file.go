package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists checks if the file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)

	// if there is an error, or the filepath doesnt exist, return false
	if err != nil || os.IsNotExist(err) {
		return false
	}
	// return false for directories - its not a file
	if info.IsDir() {
		return false
	}

	return true
}

// DirExists checks if the path exists and is a directory
func DirExists(path string) (exists bool) {
	exists = false
	info, err := os.Stat(path)

	if err == nil && info.IsDir() {
		exists = true
	}
	return
}

// FileList returns all files in a directory that have the matching extension.
//
// If `ext` is empty, then all files are returned
func FileList(directory string, ext string) (files []string) {
	files = []string{}
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil &&
			!info.IsDir() &&
			(ext == "" || filepath.Ext(info.Name()) == ext) {
			files = append(files, path)
		}
		return nil
	})
	return
}

// FileCopy copies the content of the source io.Reader the new destionation file path
//
// If the destination path is a directory or a file that already exists, this will fail
func FileCopy(source io.Reader, destinationPath string) (err error) {
	var directory string = filepath.Dir(destinationPath)
	var destination *os.File
	// if destination file exists, fail
	if FileExists(destinationPath) {
		return fmt.Errorf("destination [%s] file already exists - delete before overwriting.", destinationPath)
	}
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
