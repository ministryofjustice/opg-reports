package dbselects

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

type mockRow struct {
	ID   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type empty struct{}
type namefilter struct {
	Name string `json:"name" db:"name"`
}
type infilter struct {
	Names []string `json:"names" db:"names"`
}
type infilterx struct {
	ID    int      `json:"id" db:"id"`
	Names []string `json:"names" db:"names"`
}

const (
	mockCreate            string = `CREATE TABLE IF NOT EXISTS mocktable (id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;`
	mockSelectStmt        string = `SELECT id, name FROM mocktable ORDER BY id ASC`
	mockSelectOnNameStmt  string = `SELECT id, name FROM mocktable WHERE name = :name ORDER BY id ASC`
	mockSelectInStmt      string = `SELECT id, name FROM mocktable WHERE name IN (:names) ORDER BY id ASC`
	mockSelectInAndIDStmt string = `SELECT id, name FROM mocktable WHERE name IN (:names) AND id <= :id ORDER BY id ASC`
	mockInsertStmt        string = `INSERT INTO mocktable (name) VALUES (:name) RETURNING id`
)

var mockInserts []*dbstmts.Insert[*mockRow, int] = []*dbstmts.Insert[*mockRow, int]{
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-a"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-b"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-c"}},
	{Statement: mockInsertStmt, Data: &mockRow{Name: "mock-d"}},
}

func TestDBDBSelectIn(t *testing.T) {
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

	// -- test select that returns everything
	all := &dbstmts.Select[*empty, *mockRow]{
		Statement: mockSelectStmt,
		Data:      &empty{},
	}
	err = Select(ctx, log, db, all)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}
	if len(all.Returned) != len(mockInserts) {
		t.Errorf("select all failed, result count mismatch.")
	}

	// -- test a select with a single filter
	singleFilter := &dbstmts.Select[*namefilter, *mockRow]{
		Statement: mockSelectOnNameStmt,
		Data:      &namefilter{Name: "mock-d"},
	}
	err = Select(ctx, log, db, singleFilter)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}
	if len(singleFilter.Returned) != 1 {
		t.Errorf("select on name failed, result count mismatch.")
	}

	// -- test a select with a single filter with 0 results
	singleFilterZero := &dbstmts.Select[*namefilter, *mockRow]{
		Statement: mockSelectOnNameStmt,
		Data:      &namefilter{Name: "mock-01"},
	}
	err = Select(ctx, log, db, singleFilterZero)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}
	if len(singleFilterZero.Returned) != 0 {
		t.Errorf("select on fake name failed, result count mismatch.")
	}
	// -- test a select with an IN filter which should return 1 result
	inFilter := &dbstmts.Select[*infilter, *mockRow]{
		Statement: mockSelectInStmt,
		Data:      &infilter{Names: []string{"a", "mock-a"}},
	}
	err = Select(ctx, log, db, inFilter)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}
	if len(inFilter.Returned) != 1 {
		t.Errorf("select with in failed, result count mismatch.")
	}

	// -- test a select with an IN filter which should return 1 result
	inIDFilter := &dbstmts.Select[*infilterx, *mockRow]{
		Statement: mockSelectInAndIDStmt,
		Data:      &infilterx{Names: []string{"mock-a", "mock-d"}, ID: 2},
	}
	err = Select(ctx, log, db, inIDFilter)
	if err != nil {
		t.Errorf("unexpected select issue:\n%v", err.Error())
	}
	if len(inIDFilter.Returned) != 1 {
		t.Errorf("select with in and id limit failed, result count mismatch.")
		debugger.Dump(inIDFilter.Returned)
	}

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
	var mockSelect *dbstmts.Select[*empty, *mockRow] = &dbstmts.Select[*empty, *mockRow]{
		Statement: mockSelectStmt,
		Data:      &empty{},
	}

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
