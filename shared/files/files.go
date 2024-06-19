package files

import (
	"io/fs"
	"os"
	"path/filepath"
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
		path := filepath.Join(dir, filename)
		// fmt.Printf("saving: %s\n", path)
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
