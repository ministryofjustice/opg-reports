package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/importers/sqlite/lib"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/sources/costs"
)

func TestImportsISqliteRun(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var testFile = filepath.Join(dir, "test.json")
	var args = &lib.Arguments{}

	// -- working
	// write some dummy data to the file
	fakes := exfaker.Many[*costs.Cost](10)
	content, _ := json.MarshalIndent(fakes, "", "  ")
	os.WriteFile(testFile, content, os.ModePerm)
	// setup the args
	args = &lib.Arguments{Type: "costs", Database: dbFile, Directory: dir}
	err = Run(args)
	if err != nil {
		t.Errorf("error when running the command: [%s]", err.Error())
	}

	// -- fail base on type
	args = &lib.Arguments{Type: "unknown", Database: dbFile, Directory: dir}
	err = Run(args)
	if err == nil {
		t.Errorf("should have returned an error")
	}

	// -- fail based on fake directory
	args = &lib.Arguments{Type: "unknown", Database: dbFile, Directory: "/not-real"}
	err = Run(args)
	if err == nil {
		t.Errorf("should have returned an error")
	}
}
