package crud_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
)

type testCreate struct {
	ID   int    `json:"id,omitempty" db:"id"` // ID is a generated primary key
	Ts   string `json:"ts,omitempty" db:"ts"` // TS is timestamp when the record was created
	Name string `json:"name,omitempty" db:"name"`
}

func (self *testCreate) TableName() string {
	return "test_model"
}

func (self *testCreate) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "ts": "TEXT NOT NULL", "name": "TEXT NOT NULL"}
}
func (self *testCreate) Indexes() map[string][]string {
	return map[string][]string{
		"idx_name": {"name"},
	}
}

func TestCrudCreateTable(t *testing.T) {

	var (
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	_, err = crud.CreateTable(ctx, adaptor, &testCreate{})
	if err != nil {
		t.Errorf("unexpected error for create: [%s]", err.Error())
	}

}

func TestCrudCreateIndexes(t *testing.T) {
	var (
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	_, err = crud.CreateTable(ctx, adaptor, &testCreate{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}

	_, err = crud.CreateIndexes(ctx, adaptor, &testCreate{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

}
