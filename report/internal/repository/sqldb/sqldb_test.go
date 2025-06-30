package sqldb

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// tModel just used to test repo functions
type tModel struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func TestRepositoryNew(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test.db")
	lg := utils.Logger("ERROR", "TEXT")

	_, err = New[*tModel](ctx, lg, cfg)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}

func TestRepositoryInsertAndSelectWithTestTable(t *testing.T) {
	var (
		err     error
		bounded []*BoundStatement
		dir     = t.TempDir()
		ctx     = t.Context()
		cfg     = config.NewConfig()
		lg      = utils.Logger("ERROR", "TEXT")
	)

	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "testinsert.db")

	repo, err := New[*tModel](ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// Create
	newTbl := `CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE) STRICT;`
	_, err = repo.Exec(newTbl)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// Insert
	stmt := `INSERT INTO test (name) VALUES (:name) ON CONFLICT (name) DO UPDATE SET name=excluded.name RETURNING id;`
	bounded = []*BoundStatement{
		{Statement: stmt, Data: &tModel{Name: "test01"}},
		{Statement: stmt, Data: &tModel{Name: "test02"}},
	}

	err = repo.Insert(bounded...)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// returned values should also be set
	for _, b := range bounded {
		if b.Returned.(int64) <= 0 {
			t.Errorf("return value missing from insert")
		}
	}

	// Select without any params
	stmt = `SELECT id, name FROM test`
	sel := &BoundStatement{Statement: stmt}
	err = repo.Select(sel)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	res := sel.Returned.([]*tModel)
	if len(res) < len(bounded) {
		t.Errorf("failed to return all records")
	}

}
