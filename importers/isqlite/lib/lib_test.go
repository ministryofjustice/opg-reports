package lib_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/importers/isqlite/lib"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/sources/releases"
)

// TestImportsISqliteValidateArgs checks that the args struct is validating as expected
// with valid and invalid details
func TestImportsISqliteValidateArgs(t *testing.T) {
	var err error
	var args = &lib.Arguments{}

	// valid
	args = &lib.Arguments{Type: "costs", Database: "./database/{type}.db", Directory: "./source"}
	err = lib.ValidateArgs(args)
	if err != nil {
		t.Errorf("validate returned unexpected error: [%s]", err.Error())
	}

	args = &lib.Arguments{Type: "no-allowed", Database: "./database/{type}.db", Directory: "./source"}
	err = lib.ValidateArgs(args)
	if err == nil {
		t.Errorf("should have returned an error, but did not")
	}
}

// TestImportsISqliteGetDatabase checks that the database
// creation returns correctly and handles error conditions
func TestImportsISqliteGetDatabase(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var args = &lib.Arguments{}
	var ctx = context.Background()
	var db *sqlx.DB

	// working
	args = &lib.Arguments{Type: "costs", Database: dbFile, Directory: "./source"}
	db, err = lib.GetDatabase(ctx, args)
	if err != nil {
		t.Errorf("error with GetDatabase: [%s]", err.Error())
	}
	if db != nil {
		db.Close()
	}

	// no database path
	args = &lib.Arguments{Type: "costs", Database: "", Directory: "./source"}
	db, err = lib.GetDatabase(ctx, args)
	if err == nil {
		t.Errorf("should have generated an error about database files")
	}
	if db != nil {
		db.Close()
	}

	// no database path
	args = &lib.Arguments{Type: "not-allowed", Database: dbFile, Directory: "./source"}
	db, err = lib.GetDatabase(ctx, args)
	if err == nil {
		t.Errorf("should have generated an error about incorrect type")
	}
	if db != nil {
		db.Close()
	}

}

// TestImportsISqliteProcessDataFile test that the
// dummy generated file is inserted into the test
// database correctly and the row counts match
func TestImportsISqliteProcessDataFile(t *testing.T) {
	var err error
	var db *sqlx.DB
	var count int
	var dir = "./" //t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var testFile = filepath.Join(dir, "test.json")
	var args = &lib.Arguments{}
	var ctx = context.Background()
	var n = 1000

	// -- working
	// write some dummy data to the file
	fakes := exfaker.Many[*releases.Release](n)
	content, _ := json.MarshalIndent(fakes, "", "  ")
	os.WriteFile(testFile, content, os.ModePerm)
	// setup the args
	args = &lib.Arguments{Type: "releases", Database: dbFile, Directory: dir}
	db, err = lib.GetDatabase(ctx, args)
	if err != nil {
		t.Errorf("error getting database: [%s]", err.Error())
	}
	defer db.Close()

	// make the call to process the file
	count, err = lib.ProcessDataFile(ctx, db, args, testFile)
	if err != nil {
		t.Errorf("error processing file: [%s]", err.Error())
	}
	if count != len(fakes) {
		t.Errorf("incorrect insert count:  expected [%d] actual [%v]", len(fakes), count)
	}

}
