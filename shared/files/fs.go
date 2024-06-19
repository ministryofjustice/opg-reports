package files

import "io/fs"

type Reader interface {
	fs.FS
	fs.ReadFileFS
}

type ReaderWithDir interface {
	Reader
	BaseDir() string
}

type ReadFS struct {
	Reader
	dir string
}

func (f *ReadFS) BaseDir() string {
	return f.dir
}

func NewFS(f Reader, dir string) *ReadFS {
	return &ReadFS{
		Reader: f,
		dir:    dir,
	}
}
