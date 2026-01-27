package dbinserts

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/jmoiron/sqlx"
)

const mockTableCreate string = `
CREATE TABLE IF NOT EXISTS test_table (
	id INTEGER PRIMARY KEY,
	created_at TEXT NOT NULL DEFAULT (strftime('%FT%TZ', 'now') ),
	name TEXT NOT NULL
) STRICT;
`
const mockInsert string = `
INSERT INTO test_table (
	name
) VALUES (
	:name
)
RETURNING id;
`

type mockRow struct {
	ID        int    `json:"id,omitempty" db:"id"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	Name      string `json:"name,omitempty" db:"name"`
}

func TestDBDBInsertsWorking(t *testing.T) {
	var (
		err     error
		db      *sqlx.DB
		dir     string          = t.TempDir()
		ctx     context.Context = t.Context()
		log     *slog.Logger    = logger.New("debug", "text")
		driver  string          = "sqlite3"
		connStr string          = fmt.Sprintf("%s/%s", dir, "insert-working.db")
		mock    *mockRow        = &mockRow{Name: "test-name"}
		inserts                 = []*dbstatements.InsertStatement[*mockRow, int]{
			{Statement: mockInsert, Data: mock},
		}
	)
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()
	// db schema setup
	_, err = dbexec.Exec(ctx, log, db, dbstatements.Statement(mockTableCreate))
	if err != nil {
		t.Errorf("unexpected exec issue:\n%v", err.Error())
	}
	// now an insert
	err = Insert(ctx, log, db, inserts...)
	if err != nil {
		t.Errorf("unexpected exec issue:\n%v", err.Error())
	}
	// check the result
	if inserts[0].Returned <= 0 {
		t.Errorf("expected id to be set and positive integer")
	}
}
