package crud_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/stopwatch"
)

type testJ struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"` // ID is a generated primary key
	Name string `json:"name,omitempty" db:"name" faker:"word"`
}

type testInsert struct {
	ID   int      `json:"id,omitempty" db:"id" faker:"-"`           // ID is a generated primary key
	Ts   string   `json:"ts,omitempty" db:"ts" faker:"time_string"` // TS is timestamp when the record was created
	Name string   `json:"name,omitempty" db:"name" faker:"word"`
	List []*testJ `json:"list,omitempty" db:"-" faker:"slice_len=1"`
}

func (self *testInsert) UniqueValue() string {
	return self.Name
}
func (self *testInsert) UniqueField() string {
	return "name"
}

func (self *testInsert) UpsertUpdate() string {
	return "name=excluded.name"
}

func (self *testInsert) TableName() string {
	return "test_model"
}

func (self *testInsert) GetID() int {
	return self.ID
}
func (self *testInsert) SetID(id int) {
	self.ID = id
}
func (self *testInsert) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "ts": "TEXT NOT NULL", "name": "TEXT NOT NULL UNIQUE"}
}
func (self *testInsert) Indexes() map[string][]string {
	return map[string][]string{
		"idx_name": {"name"},
	}
}
func (self *testInsert) InsertColumns() []string {
	return []string{
		"name", "ts",
	}
}
func (self *testInsert) New() dbs.Cloneable {
	return &testInsert{}
}

var _ dbs.InsertableRow = &testInsert{}
var _ dbs.Cloneable = &testInsert{}
var _ dbs.CreateableTable = &testInsert{}

// TestCrudInsertMany checks testing 500,000 records being created
// to check performance record errors / timeouts
func TestCrudInsertMany(t *testing.T) {
	fakerextras.AddProviders()
	var (
		n       int = 500000
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		tests   []*testInsert   = fakermany.Fake[*testInsert](n)
		results []*testInsert
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	_, err = crud.CreateTable(ctx, adaptor, &testInsert{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testInsert{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	stopwatch.Start()
	results, err = crud.Insert(ctx, adaptor, &testInsert{}, tests...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}
	// should be about 4.5 seconds..
	fmt.Println(stopwatch.Seconds())

	if len(tests) != len(results) {
		t.Errorf("differing number of result to tests - expected [%d] actual [%v]", len(tests), len(results))
	}

}

func TestCrudInsertSimple(t *testing.T) {
	var (
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		tests   []*testInsert   = []*testInsert{
			{Name: "testA", Ts: time.Now().Format(time.RFC3339)},
			{Name: "testB", Ts: time.Now().Format(time.RFC3339)},
		}
		results []*testInsert
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	_, err = crud.CreateTable(ctx, adaptor, &testInsert{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testInsert{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &testInsert{}, tests...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(tests) != len(results) {
		t.Errorf("differing number of result to tests - expected [%d] actual [%v]", len(tests), len(results))
	}

	for _, expected := range tests {
		var found = false
		for _, actual := range results {
			if actual.Name == expected.Name {
				found = true
			}
		}
		if !found {
			t.Errorf("test was not found in the results")
		}
	}

	for _, res := range results {
		if res.ID <= 0 {
			t.Errorf("database id not set on record")
		}
	}
}
