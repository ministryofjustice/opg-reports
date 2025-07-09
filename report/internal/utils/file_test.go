package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opg-reports/report/internal/utils"
)

func TestFileExists(t *testing.T) {
	var dir = t.TempDir()

	// make sure a file that doesnt exist is not found
	doesntExist := filepath.Join(dir, "not-a-real-file.txt")
	if utils.FileExists(doesntExist) {
		t.Errorf("found a file that should not exist")
	}

	// make sure a file that does exist is found
	doesExist := filepath.Join(dir, "file-exists.txt")
	os.WriteFile(doesExist, []byte("test file."), os.ModePerm)
	if !utils.FileExists(doesExist) {
		t.Errorf("failed to find a file that should exist")
	}

	// make sure it returns false for a directory
	subDir := filepath.Join(dir, "sub-dir")
	os.MkdirAll(subDir, os.ModePerm)
	if utils.FileExists(subDir) {
		t.Errorf("should return false for a directory")
	}

}

func TestFileCopy(t *testing.T) {

	var (
		err     error
		dest    string
		dir     = t.TempDir()
		content = "test file."
		src     = filepath.Join(dir, "source.txt")
	)
	// create source file and pointer
	os.WriteFile(src, []byte(content), os.ModePerm)
	f, _ := os.Open(src)
	defer f.Close()

	// copy to a valid location should work
	dest = filepath.Join(dir, "t1.txt")
	err = utils.FileCopy(f, dest)
	if err != nil {
		t.Errorf("unexpected error on file copy: %s", err.Error())
	}

	// make sure copy contains same content
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Errorf("unexpected error on reading copy: %s", err.Error())
	}
	if string(b) != content {
		t.Errorf("expected copy content to match exactly: src [%s] dest [%s]", content, string(b))
	}

	// create a file with content in set location and try to overwrite - this should throw
	// an error
	dest = filepath.Join(dir, "t2.txt")
	os.WriteFile(dest, []byte("test2"), os.ModePerm)
	err = utils.FileCopy(f, dest)
	if err == nil {
		t.Errorf("expected an error when trying to overwrite a file")
	}

	// create a directory and try to copy to that, this should fail
	dest = filepath.Join(dir, "sub")
	os.MkdirAll(dest, os.ModePerm)
	err = utils.FileCopy(f, dest)
	if err == nil {
		t.Errorf("expected an error when copying to a directory path")
	}
	if err != nil && !strings.Contains(err.Error(), "directory") {
		t.Errorf("expected error to relate to directory: %s", err.Error())
	}

}
