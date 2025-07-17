package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"opg-reports/report/internal/utils"
)

type testTar struct {
	Name string `json:"name"`
}

func TestTarGzSuccess(t *testing.T) {
	var (
		err         error
		dir         = t.TempDir()
		srcDir      = filepath.Join(dir, "src")
		srcFileA    = filepath.Join(srcDir, "a.json")
		srcFileB    = filepath.Join(srcDir, "b.json")
		srcDataA    = []*testTar{{Name: "A"}, {Name: "B"}}
		srcDataB    = []*testTar{{Name: "A"}, {Name: "B"}, {Name: "C"}}
		destArchive = filepath.Join(dir, "archive", "archive.tar.gz")
		extractTo   = filepath.Join(dir, "extract")
	)

	// generate a dummy files with some json in
	utils.StructToJsonFile(srcFileA, srcDataA)
	utils.StructToJsonFile(srcFileB, srcDataB)

	// now compress those files into a tar ball
	err = utils.TarGzCreate(destArchive, []string{srcFileA, srcFileB})
	if err != nil {
		t.Errorf("unexpected error compressing files: %s", err.Error())
	}
	if !utils.FileExists(destArchive) {
		t.Errorf("archive file was not created in expected location")
	}
	// open archive
	f, err := os.Open(destArchive)
	if err != nil {
		t.Errorf("unexpected error opening archive: %s", err.Error())
	}
	defer f.Close()
	// now extract to a new folder
	err = utils.TarGzExtract(extractTo, f)
	if err != nil {
		t.Errorf("unexpected error extracting files: %s", err.Error())
	}

	// check files extracted in know places
	extractA := filepath.Join(dir, "extract", srcFileA)
	extractB := filepath.Join(dir, "extract", srcFileB)
	if !utils.FileExists(extractA) {
		t.Errorf("could not find extracted file: %s", extractA)
	}
	if !utils.FileExists(extractB) {
		t.Errorf("could not find extracted file: %s", extractB)
	}

	// now load up files and compare to original content to make sure
	// number of structs match
	structA := []*testTar{}
	err = utils.StructFromJsonFile(extractA, &structA)
	if err != nil {
		t.Errorf("unexpected error compressing files: %s", err.Error())
	}
	if len(structA) != len(srcDataA) {
		t.Errorf("number of records in extracted file did not match")
	}
	structB := []*testTar{}
	err = utils.StructFromJsonFile(extractB, &structB)
	if err != nil {
		t.Errorf("unexpected error compressing files: %s", err.Error())
	}
	if len(structB) != len(srcDataB) {
		t.Errorf("number of records in extracted file did not match")
	}

}
