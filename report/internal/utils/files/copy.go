package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Copy(source io.Reader, destinationPath string) (err error) {
	var directory string = filepath.Dir(destinationPath)
	var destination *os.File
	// if destination file exists, fail
	if Exists(destinationPath) {
		err = fmt.Errorf("destination [%s] file already exists - delete before overwriting.", destinationPath)
		err = errors.Join(ErrFileExists, err)
		return
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
