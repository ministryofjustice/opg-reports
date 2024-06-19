package files

import "io/fs"

// IReadFS (Interface File System) is a wrpping interface to ensure
// the filesystem being used impliments key functions from
// fs.FS & fs.ReadFileFS
type IReadFS interface {
	fs.FS
	fs.ReadFileFS
}

// IWriteFS extends on IReadFS with a directory function to enabled
// writing data out to set paths
type IWriteFS interface {
	IReadFS
	BaseDir() string
}

type ReadFS struct {
	IReadFS
	dir string
}

func (f *ReadFS) BaseDir() string {
	return f.dir
}

func NewFS(f IReadFS, dir string) *ReadFS {
	return &ReadFS{
		IReadFS: f,
		dir:     dir,
	}
}
