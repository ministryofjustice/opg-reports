// Package files provides method to return all files, or fitler list of files from a T IReadFS filesystem.
//
// As the api and front will swap between file system based on directory or an embedded version, this package
// contains interfaces (IReadFS, IWriteFS) and concrete versions (WriteFS) that both embed.FS and result of
// os.DirFs can be casted into to simplify that issue.
//
// The functions within this package utilise those interfaces so and are consumed by the api and data store
// objects
package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// All returns all files recursively from the base of the filesystem
// Setting jsonOnly to true will limit the returned files to those
// with a .json extension
func All[T IReadFS](f T, jsonOnly bool) (files []*PathFile) {
	var jsonExt string = ".json"

	files = []*PathFile{}

	fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		// no errors & is a file
		if err == nil && !d.IsDir() {
			// jsonOnly is false, or its true & its a json file
			if !jsonOnly || (jsonOnly && filepath.Ext(d.Name()) == jsonExt) {
				files = append(files, NewFile(d, path))
			}
		}
		return nil
	})

	return files
}

// Filter uses a list of files passed and returns those that match *all* filterPatterns
// that are passed along
func Filter(all []*PathFile, filterPatterns ...string) []*PathFile {
	fSet := []*PathFile{}

	// a file can only be added to the final set if it matches all filter patterns passed along
	for _, file := range all {
		add := false
		matchesAll := true

		for _, pattern := range filterPatterns {
			add = true
			reg, err := regexp.Compile(pattern)
			// if theres an error, or this is file is a dir, or it doesnt match the pattern
			// then it should not be added
			if err != nil || file.IsDir() || !reg.MatchString(file.Name()) {
				matchesAll = false
			}
		}
		if add && matchesAll {
			fSet = append(fSet, file)
		}

	}
	return fSet

}

// ReadFile wraps the ReadFS ReadFile
func ReadFile[T IReadFS](fS T, file *PathFile) ([]byte, error) {
	return fS.ReadFile(file.Path)
}

// WriteFile will create the directory path, join the dir and filename together
// and then try to use os.WriteFile to write content to that file
func WriteFile(dir string, filename string, content []byte) (err error) {

	if e := os.MkdirAll(dir, os.ModePerm); e != nil {
		err = e
	} else {

		filename = strings.ReplaceAll(filename, dir, "")
		path := filepath.Join(dir, filename)
		err = os.WriteFile(path, content, os.ModePerm)
	}

	return
}

// SaveFile presumes the directory path described in f.Path exists and is
// writable. Wraps os.WriteFile otherwise
func SaveFile[T IWriteFS](fS T, f *PathFile, content []byte) (err error) {
	dir := fS.BaseDir()
	name := strings.ReplaceAll(f.Path, dir, "")
	return WriteFile(dir, name, content)
}
