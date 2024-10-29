package datastore_test

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
)

const (
	dummyDbCreate datastore.CreateStatement      = `CREATE TABLE IF NOT EXISTS dummy (id INTEGER PRIMARY KEY, label TEXT NOT NULL) STRICT;`
	dummyDbInsert datastore.InsertStatement      = `INSERT INTO dummy(label) VALUES(:label) RETURNING id`
	dummyGetOne   datastore.SelectStatement      = `SELECT count(*) FROM dummy ORDER BY id ASC LIMIT 1`
	dummyGet      datastore.NamedSelectStatement = `SELECT id, label FROM dummy WHERE id > :min ORDER BY id ASC`
)

type dummy struct {
	ID    int    `json:"id,omitempty" db:"id"`
	Label string `json:"label,omitempty" db:"label" faker:"word"`
}

type dummyE struct {
	ID           int    `json:"id,omitempty" db:"id"`
	Label        string `json:"label,omitempty" db:"label" faker:"word"`
	Organisation string `json:"organisation,omitempty"`
	Unit         string `json:"unit,omitempty"`
}

type dummyP struct {
	Min int `json:"min" db:"min"`
}

type testStruct struct {
	StartDate  string `json:"start_date,omitempty" db:"start_date"`
	EndDate    string `json:"end_date,omitempty" db:"end_date"`
	DateFormat string `json:"date_format,omitempty" db:"date_format"`
}

// TestDatastoreDB runs a series of simple create,
// insert and select queries to ensure a database
// is setup as expected
func TestDatastoreDB(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "test.db")
	var ctx context.Context = context.Background()
	// -- single check
	var record *dummy = &dummy{Label: "One"}
	// -- multiple check
	var n int = 15000
	var records []*dummy = []*dummy{}

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("unexpected error creating db in datastore: [%s]", err.Error())
	}

	datastore.Create(ctx, db, []datastore.CreateStatement{dummyDbCreate})

	// -- insert one
	id, err := datastore.InsertOne(ctx, db, dummyDbInsert, record, nil)
	if err != nil {
		t.Errorf("error inserting: [%s]", err.Error())
	}
	if id != 1 {
		t.Errorf("expected first insert to have id of 1, actual [%d]", id)
	}

	// -- insert many
	records = exfaker.Many[dummy](n)
	ids, err := datastore.InsertMany(ctx, db, dummyDbInsert, records)
	if err != nil {
		t.Errorf("error inserting many [%s]", err.Error())
	}
	if len(ids) != n {
		t.Errorf("did not insert all rows - expected [%d] actual [%d]", n, len(ids))
	}

	// -- get one, no args
	got, err := datastore.Get[int](ctx, db, dummyGetOne)
	if err != nil {
		t.Errorf("error with get [%s]", err.Error())
	}
	if got != n+1 {
		t.Errorf("incorrect db count")
	}
	// -- get many with named params
	p := &dummyP{Min: 1}
	found, err := datastore.Select[[]*dummy](ctx, db, dummyGet, p)
	if err != nil {
		t.Errorf("error with select [%s]", err.Error())
	}
	if len(found) != got-1 {
		t.Errorf("select failed to get everything")
	}

}

// TestDatastoreDBNeeds checks that named values within
// a NamedSelectStatement are matched as expected including
// odd naming patterns
func TestDatastoreDBNeeds(t *testing.T) {

	var testString datastore.NamedSelectStatement = ""
	var found []string = []string{}

	// this should return min & max
	testString = `SELECT count(*) FROM table WHERE id > :min AND id < :max`
	found = datastore.Needs(testString)
	if !slices.Contains(found, "min") || !slices.Contains(found, "max") {
		t.Errorf("needs did not find all fields: [%s]", strings.Join(found, ","))
	}
	if len(found) != 2 {
		t.Errorf("incorrect number of fields found: [%s]", strings.Join(found, ","))
	}
	// testing multiline string
	testString = `
	SELECT
		count(*)
	FROM table
	WHERE
		id > :min
		AND id < :max
		AND date_end is TRUE
		AND start > :start_dateWith-oddName
	`
	found = datastore.Needs(testString)
	if !slices.Contains(found, "min") || !slices.Contains(found, "max") || !slices.Contains(found, "start_dateWith-oddName") {
		t.Errorf("needs did not find all fields: [%s]", strings.Join(found, ","))
	}
	if len(found) != 3 {
		t.Errorf("incorrect number of fields found: [%s]", strings.Join(found, ","))
	}

}

// TestDatastoreDBValidateParameters checks that valid
// param and needs as well as invalid versions are triggered
// correctly
func TestDatastoreDBValidateParameters(t *testing.T) {
	var err error
	var needs []string = []string{}
	var params *testStruct = &testStruct{}

	// -- test it works and ignores extra fields
	needs = []string{"start_date"}
	params = &testStruct{StartDate: "test", EndDate: "test"}
	err = datastore.ValidateParameters(params, needs)
	if err != nil {
		t.Errorf("param should be valid: [%s]", err.Error())
	}

	// -- test a failing one
	needs = []string{"end_date"}
	params = &testStruct{StartDate: "test"}
	err = datastore.ValidateParameters(params, needs)
	if err == nil {
		t.Errorf("param should throw error, but didnt")
	}

}

// TestDatastoreColumnValues checks that all permutations of column
// values are found correctly for the data
func TestDatastoreColumnValues(t *testing.T) {

	simple := []*dummyE{
		{Organisation: "A", Unit: "A", Label: "One"},
		{Organisation: "A", Unit: "B", Label: "One"},
		{Organisation: "A", Unit: "C", Label: "One"},
		{Organisation: "A", Unit: "1", Label: "Three"},
		{Organisation: "A", Unit: "Z", Label: "One"},
		{Organisation: "A", Unit: "Z", Label: "Four"},
	}

	cols := []string{"organisation", "unit", "label"}

	actual := datastore.ColumnValues(simple, cols)
	expected := map[string]int{"organisation": 1, "unit": 5, "label": 3}
	for col, count := range expected {
		l := len(actual[col])
		if l != count {
			t.Errorf("[%s] values dont match - expected [%d] actual [%v]", col, count, l)
			fmt.Println(actual[col])
		}
	}

}
