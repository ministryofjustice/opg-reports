package files

import "io/fs"

type PathFile struct {
	fs.DirEntry
	Path string
}

func NewFile(f fs.DirEntry, path string) *PathFile {
	return &PathFile{
		DirEntry: f,
		Path:     path,
	}
}
