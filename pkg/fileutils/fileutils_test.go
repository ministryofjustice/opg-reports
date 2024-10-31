package fileutils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/pkg/fileutils"
)

// TestFileUtilsCopy checks the copy function does make an exact match of
// an existing file
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

// TestFileUtilsDownloadFromUrl downloads a known url from govuk
// website to a local file and checks its present and has content
func TestFileUtilsDownloadFromUrl(t *testing.T) {
	var err error
	var path string
	var b []byte
	var content string
	var testUrl = "https://www.gov.uk/power-of-attorney/make-lasting-power"
	var dir = t.TempDir()
	var name = "make.html"

	path, err = fileutils.DownloadFromUrl(testUrl, dir, name, time.Second*2)
	if err != nil {
		t.Errorf("error downloading file")
	}

	if !fileutils.Exists(path) {
		t.Errorf("file was not downloaded")
	}

	b, err = os.ReadFile(path)
	if err != nil {
		t.Errorf("error reading file")
	}
	content = string(b)
	if len(content) <= 0 {
		t.Errorf("error with file content")
	}
}

// TestFileUtilsZipCreateExtract generates a test zip file
// from a series of files and directories and then
// extracts that zip and compares the results to ensure
// they match with original versions
//   - empty directories are tested, but very simple
func TestFileUtilsZipCreateExtract(t *testing.T) {
	var err error
	// td := os.TempDir()
	// dir2, _ := os.MkdirTemp(td, "test-*")

	var dir2 = t.TempDir()
	var dir = t.TempDir()
	var zip = filepath.Join(dir, "test.zip")
	var zip2 = filepath.Join(dir2, "openme.zip")
	var content = "foobar"
	var extracted []string

	var files = []string{
		filepath.Join(dir, "tmp.txt"),
		filepath.Join(dir, "test.csv"),
		filepath.Join(dir, "sub/file.html"),
		filepath.Join(dir, "sub/sub/tmp.json"),
		filepath.Join(dir, "empty-directory"),
		filepath.Join(dir, "empty/sub-directory"),
	}
	// create the file in the test directory
	for _, file := range files {
		var isFile = strings.Contains(file, ".")
		if isFile {
			os.MkdirAll(filepath.Dir(file), os.ModePerm)
			os.WriteFile(file, []byte(content), os.ModePerm)
		} else {
			os.MkdirAll(file, os.ModePerm)
		}
	}

	err = fileutils.ZipCreate(zip, files, dir)
	if err != nil {
		t.Errorf("error making zip: [%s]", err.Error())
	}

	if !fileutils.Exists(zip) {
		t.Errorf("zip was not created")
	}
	// copy to a new directory
	fileutils.CopyFromPath(zip, zip2)
	// Now we extract the zip to compare the files
	extracted, err = fileutils.ZipExtract(zip2, dir2+"/extracted/")
	if err != nil {
		t.Errorf("error extracting zip: [%s]", err.Error())
	}

	if len(extracted) != len(files) {
		t.Errorf("extracted files mis match")
	}

	// now we check each file / directory exists with correct content in the destination
	for _, file := range files {
		var originallyDir = fileutils.IsDir(file)
		var newPath = strings.ReplaceAll(file, dir, dir2+"/extracted")
		var nowDir = fileutils.IsDir(newPath)

		// fmt.Printf("[%v => %v] %s \n", originallyDir, nowDir, newPath)

		if originallyDir != nowDir {
			t.Errorf("error on determining directory or not [%s] => [%s]", file, newPath)
		}
		// if its not a directory, check the content
		if !nowDir {
			ct, err := os.ReadFile(newPath)
			if err != nil {
				t.Errorf("error reading file [%s]: [%s]", newPath, err.Error())
			}
			if string(ct) != content {
				t.Errorf("file content mismatch: [%s]", newPath)
			}
		}

	}
}
