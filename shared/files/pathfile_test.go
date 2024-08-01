package files

import (
	"io/fs"
	"opg-reports/shared/logger"
	"os"
	"testing"
)

func TestNewPathFile(t *testing.T) {
	logger.LogSetup()
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
