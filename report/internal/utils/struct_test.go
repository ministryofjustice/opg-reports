package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type testStruct struct {
	Name string `json:"name,omitempty"`
}

func TestStructToJsonFile(t *testing.T) {
	var (
		err      error
		file     string
		dir      = t.TempDir()
		itemIn   = &testStruct{Name: "test1"}
		itemOut  = &testStruct{}
		itemsIn  = []*testStruct{{Name: "test2"}, {Name: "test3"}}
		itemsOut = []*testStruct{}
	)

	// test single item being written to a file
	file = filepath.Join(dir, "single.json")
	err = utils.StructToJsonFile(file, itemIn)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !utils.FileExists(file) {
		t.Errorf("file was not created [%s]", file)
	}
	// read in the file to compare
	err = utils.StructFromJsonFile(file, &itemOut)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if itemOut.Name != itemIn.Name {
		t.Errorf("original and imported dont match: expected [%s] actual [%s]", itemIn, itemOut)
	}

	// test multiple items being written to a file
	file = filepath.Join(dir, "many.json")
	err = utils.StructToJsonFile(file, itemsIn)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !utils.FileExists(file) {
		t.Errorf("file was not created [%s]", file)
	}
	// read in the file to compare
	err = utils.StructFromJsonFile(file, &itemsOut)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(itemsOut) != len(itemsIn) {
		t.Errorf("original and imported dont match:")
	}

}
