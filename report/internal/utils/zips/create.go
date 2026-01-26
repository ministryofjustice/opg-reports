package zips

import (
	"archive/zip"
	"fmt"
	"io"
	"opg-reports/report/internal/utils/files"
	"os"
	"path/filepath"
	"strings"
)

// Create generates a zipfile from the list of files passed along.
// The filepaths need to be the fully accessible (best to go from the root)
// You can trim the filepath pushed to the zip be adding that to the removePath
func Create(zipFilepath string, filepaths []string, removePath string) (err error) {
	var (
		archive      *os.File
		zipWriter    *zip.Writer
		zipDirectory = filepath.Dir(zipFilepath)
	)

	// try to make the directory path
	if err = os.MkdirAll(zipDirectory, os.ModePerm); err != nil {
		return
	}
	// try to create the zipfile
	if archive, err = os.Create(zipFilepath); err != nil {
		return
	}
	// create the zip writer
	zipWriter = zip.NewWriter(archive)
	// close
	defer archive.Close()
	defer zipWriter.Close()

	for _, file := range filepaths {
		if err = writeToZip(zipWriter, file, removePath); err != nil {
			return
		}
	}

	return
}

// writeToZip handles writing a single file/directory into a zip file using
// the zipwriter
func writeToZip(zipWriter *zip.Writer, file string, removePath string) (err error) {
	var localFile *os.File
	var zipFile io.Writer
	var isDir bool = files.DirExists(file)
	var zfile string = strings.ReplaceAll(file, removePath, "")

	// if its a directory, append the /
	if isDir {
		zfile = fmt.Sprintf("%s%c", zfile, os.PathSeparator)
	}

	// create a stub in the zip - this is all a directory requires
	if zipFile, err = zipWriter.Create(zfile); err != nil {
		return
	}

	// if its not a directory (therefore a file), copy the contents
	// - if file opening fails, return
	// - if copy fails, return
	if !isDir {
		if localFile, err = os.Open(file); err != nil {
			return
		}
		defer localFile.Close()
		if _, err = io.Copy(zipFile, localFile); err != nil {
			return
		}
	}
	return
}
