package fileutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/fileutils"
)

// Check the copy function does make an exact match of an existing file
func TestFileUtilsCopy(t *testing.T) {
	var err error
	var sourceFile *os.File
	var dir = t.TempDir()
	var source = filepath.Join(dir, "source.txt")
	var destPath = filepath.Join(dir, "dest.txt")
	var expected = "TEST"

	err = os.WriteFile(source, []byte(expected), os.ModePerm)
	if err != nil {
		t.Errorf("error creating test file [%s]", err.Error())
	}

	sourceFile, _ = os.Open(source)
	_, err = fileutils.Copy(sourceFile, destPath)
	if err != nil {
		t.Errorf("error copying [%s]", err.Error())
	}

	b, _ := os.ReadFile(destPath)
	actual := string(b)
	if actual != expected {
		t.Errorf("error with copy - expected [%s] actual [%v]", expected, actual)
	}
}
