package crud_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/pretty"
)

type testSelect struct {
	ID   int    `json:"id,omitempty" db:"id" faker:"-"` // ID is a generated primary key
	Name string `json:"name,omitempty" db:"name" faker:"word"`
}

func (self *testSelect) TableName() string {
	return "test_selects"
}
func (self *testSelect) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "name": "TEXT NOT NULL"}
}
func (self *testSelect) Indexes() map[string][]string {
	return map[string][]string{
		"idx_name": {"name"},
	}
}
func (self *testSelect) GetID() int {
	return self.ID
}
func (self *testSelect) SetID(id int) {
	self.ID = id
}
func (self *testSelect) InsertColumns() []string {
	return []string{"name"}
}
func (self *testSelect) New() dbs.Cloneable {
	return &testSelect{}
}

type testParam struct {
	ID int `json:"id,omitempty" db:"id" faker:"-"`
}

func TestCrudSelectMultipleStructs(t *testing.T) {
	fakerextras.AddProviders()
	var (
		n       int = 100
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		tests   []*testSelect   = fakermany.Fake[*testSelect](n)
		results []*testSelect
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &testSelect{}, tests...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}
	if len(tests) != len(results) {
		t.Errorf("differing number of result to tests - expected [%d] actual [%v]", len(tests), len(results))
	}

	stmt := `SELECT * FROM test_selects WHERE ID > :id`
	results, err = crud.Select[*testSelect](ctx, adaptor, stmt, &testParam{ID: 5})

	if len(results) != n-5 {
		t.Errorf("incorrect number of results returned.")
	}

	stmt = `SELECT * FROM test_selects WHERE ID > 2`
	results, err = crud.Select[*testSelect](ctx, adaptor, stmt, nil)

	if len(results) != n-2 {
		t.Errorf("incorrect number of results returned.")
	}

}

func TestCrudSelectSingleStructs(t *testing.T) {
	fakerextras.AddProviders()
	var (
		n       int = 100
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		tests   []*testSelect   = fakermany.Fake[*testSelect](n)
		results []*testSelect
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &testSelect{}, tests...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}
	if len(tests) != len(results) {
		t.Errorf("differing number of result to tests - expected [%d] actual [%v]", len(tests), len(results))
	}

	stmt := `SELECT * FROM test_selects WHERE ID = :id`
	results, err = crud.Select[*testSelect](ctx, adaptor, stmt, &testParam{ID: 10})

	if len(results) != 1 {
		t.Errorf("expected only a single result")
		pretty.Print(results)
	}

}

func TestCrudSelectPrimatives(t *testing.T) {
	fakerextras.AddProviders()
	var (
		n       int = 100
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		tests   []*testSelect   = fakermany.Fake[*testSelect](n)
		results []*testSelect
		res     []int
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testSelect{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &testSelect{}, tests...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}
	if len(tests) != len(results) {
		t.Errorf("differing number of result to tests - expected [%d] actual [%v]", len(tests), len(results))
	}

	stmt := `SELECT COUNT(*) FROM test_selects LIMIT 1;`
	res, err = crud.Select[int](ctx, adaptor, stmt, nil)

	pretty.Print(res)

}
