package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

// writeToZip handles writing a single file/directory into a zip file using
// the zipwriter
func writeToZip(zipWriter *zip.Writer, file string, removePath string) (err error) {
	var localFile *os.File
	var zipFile io.Writer
	var isDir bool = DirExists(file)
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
		if err = FileCopy(zipFile, path); err != nil {
			return
		}
	}
	return
}
