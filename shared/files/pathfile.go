package files

import "io/fs"

// PathFile is an extended file that also tracks the original
// path it was found at seperately to make swapping between
// embedded and dir passed filesystems easier
type PathFile struct {
	fs.DirEntry
	Path string
}

// NewFile returns a instance with the dir entry and extra path info
func NewFile(f fs.DirEntry, path string) *PathFile {
	return &PathFile{
		DirEntry: f,
		Path:     path,
	}
}
