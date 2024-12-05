// Package fileutils include extra features for handling files
//
// Includes helpers like copying a buffer to a file
package fileutils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/fetch"
)

// Exists checks if the file exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDir(path string) (isDir bool) {
	isDir = false
	if f, err := os.Stat(path); err == nil {
		isDir = f.IsDir()
	}
	return
}

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

// CopyFromPath uses filepaths for both source and destination.
// Tries to open the source and the passes along to Copy
func CopyFromPath(source string, destination string) (err error) {
	var sourceFile *os.File
	// if its a directory, just create it
	if IsDir(source) {
		err = os.MkdirAll(destination, os.ModePerm)
		return
	}

	sourceFile, err = os.Open(source)
	if err != nil {
		return
	}
	_, err = Copy(sourceFile, destination)
	return
}

// DownloadFromUrl fetches the content present at the url and creates a local copy
func DownloadFromUrl(url string, destinationDir string, destinationName string, timeout time.Duration) (path string, err error) {
	var response *http.Response

	path = filepath.Join(destinationDir, destinationName)
	// get the http response
	response, err = fetch.Response(url, timeout)
	defer response.Body.Close()
	if err != nil {
		return
	}
	// copy the response body content into a new filepath
	_, err = Copy(response.Body, path)
	return
}

// ZipCreate generates a zipfile from the list of files passed along.
// The filepaths need to be the fully accessible (best to go from the root)
// You can trim the filepath pushed to the zip be adding that to the removePath
func ZipCreate(zipFilepath string, filepaths []string, removePath string) (err error) {
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
	var isDir bool = IsDir(file)
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

// ZipExtract extracts all of the files & directories within the
// zip (at zipFilepath) into the destinationDir
// If the destination does not exist it will try to create it
// Returns list of all files / directories that are extracted
func ZipExtract(zipFilepath string, destinationDir string) (extracted []string, err error) {
	var archive *zip.ReadCloser
	extracted = []string{}
	// open the zip
	archive, err = zip.OpenReader(zipFilepath)
	if err != nil {
		return
	}
	defer archive.Close()
	// create the destination directory
	if err = os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return
	}
	// extract each file at a time
	for _, file := range archive.File {
		if err = extractFromZip(file, destinationDir); err != nil {
			return
		}
		extracted = append(extracted, file.Name)
	}
	return
}

// extractFromZip handles extracting a single from the zip into the
// destination directory
func extractFromZip(file *zip.File, destinationDir string) (err error) {
	var zipFile io.ReadCloser
	var isDir = file.FileInfo().IsDir()
	var path = filepath.Join(destinationDir, file.Name)

	if isDir {
		os.MkdirAll(path, os.ModePerm)
	} else {
		// make the parent dir
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if zipFile, err = file.Open(); err != nil {
			return
		}
		if _, err = Copy(zipFile, path); err != nil {
			return
		}
	}
	return
}
