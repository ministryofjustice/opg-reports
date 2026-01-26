package zips

import (
	"archive/zip"
	"io"
	"opg-reports/report/internal/utils/files"
	"os"
	"path/filepath"
)

// Extract extracts all of the files & directories within the
// zip (at zipFilepath) into the destinationDir
// If the destination does not exist it will try to create it
// Returns list of all files / directories that are extracted
func Extract(zipFilepath string, destinationDir string) (extracted []string, err error) {
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
		if err = files.Copy(zipFile, path); err != nil {
			return
		}
	}
	return
}
