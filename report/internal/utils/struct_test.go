package utils_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"opg-reports/report/internal/utils"
)

type testStruct struct {
	Name string `json:"name,omitempty"`
}

func TestStructFromAndToJsonFiles(t *testing.T) {
	var (
		err      error
		file     string
		dir      = t.TempDir()
		itemIn   = &testStruct{Name: "test1"}
		itemOut  = &testStruct{}
		itemsIn  = []*testStruct{{Name: "test2"}, {Name: "test3"}}
		itemsOut = []*testStruct{}
	)

	// test a failing structojson - something that isnt a struct
	file = filepath.Join(dir, "fail.json")
	err = utils.StructToJsonFile(file, func() { fmt.Print("here") })
	if err == nil {
		t.Errorf("expected an error when non-struct used")
	}

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

	// test loading a file doesnt exist throws an error
	err = utils.StructFromJsonFile("foo-bar.md", &itemsOut)
	if err == nil {
		t.Errorf("expected error for a file that doesnt exist")
	}

}
