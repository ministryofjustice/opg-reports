package files

import (
	"io/fs"
	"os"
	"testing"
)

func TestNewPathFile(t *testing.T) {
	dir := "testdata"
	dirFs := os.DirFS(dir)
	entries, err := fs.ReadDir(dirFs, ".")
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	for _, f := range entries {
		fp := NewFile(f, f.Name())
		if fp.Name() != f.Name() {
			t.Errorf("filenames dont match")
		}
	}
}
