package dbselects

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

type mockRow struct {
	ID   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type empty struct{}

const (
	mockCreate     string = `CREATE TABLE IF NOT EXISTS mocktable (id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;`
	mockSelectStmt string = `SELECT id, name FROM mocktable ORDER BY id ASC`
	mockInsertStmt string = `INSERT INTO mocktable (name) VALUES (:name) RETURNING id`
)

var mockInserts []*dbstmts.Insert[*mockRow, int] = []*dbstmts.Insert[*mockRow, int]{
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-a"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-b"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-c"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-d"}},
}

var mockSelect *dbstmts.Select[*empty, *mockRow] = &dbstmts.Select[*empty, *mockRow]{
	Statement: mockSelectStmt,
	Data:      &empty{},
}

func TestDBDBSelectsWorking(t *testing.T) {

	var (
		err     error
		db      *sqlx.DB
		dir     string          = t.TempDir()
		ctx     context.Context = t.Context()
		log     *slog.Logger    = logger.New("error")
		driver  string          = "sqlite3"
		connStr string          = fmt.Sprintf("%s/%s", dir, "db-select-working.db")
	)
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()

	// run the create
	_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(mockCreate))
	if err != nil {
		t.Errorf("unexpected exec issue:\n%v", err.Error())
	}
	// run set of inserts
	err = dbinserts.Insert(ctx, log, db, mockInserts...)
	if err != nil {
		t.Errorf("unexpected insert issue:\n%v", err.Error())
	}
	// select all the things
	err = Select(ctx, log, db, mockSelect)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}

	if len(mockSelect.Returned) != len(mockInserts) {
		t.Errorf("insert & select failed with mis matching counts")
	}

}
