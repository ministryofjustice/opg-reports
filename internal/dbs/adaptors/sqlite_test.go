package adaptors

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

// interface testing
var (
	_ dbs.Adaptor   = &Sqlite{}
	_ dbs.Formatter = &SqliteFormatting{}
)

func TestAdaptorsSqlite(t *testing.T) {
	var (
		err        error
		sq         *Sqlite
		connect    string
		ctx        context.Context = context.Background()
		driver     string          = "sqlite3"
		paramChunk string          = "_journal=WAL"
		dir        string          = t.TempDir()
		path       string          = filepath.Join(dir, "test.db")
	)

	sq, err = NewSqlite(path, true)
	if err != nil {
		t.Errorf("error calling new: [%s]", err.Error())
	}
	// test it created a db file
	if !fileutils.Exists(path) {
		t.Errorf("database file was not created")
	}
	// test driver name - this is fixed on creation, so should always match
	if driver != sq.Connector().DriverName() {
		t.Errorf("drivers did not match - expected [%s] actual [%v]", driver, sq.Connector().DriverName())
	}

	// test the connection string has both path and param chunk
	connect = sq.Connector().String()
	if !strings.Contains(connect, path) {
		t.Errorf("connection did not match - expected to contains [%s] actual [%v]", path, connect)
	}
	if !strings.Contains(connect, paramChunk) {
		t.Errorf("connection did not match - expected to contains [%s] actual [%v]", paramChunk, connect)
	}
	// check the db fetch worked
	if _, err = sq.DB().Get(ctx, sq.Connector()); err != nil {
		t.Errorf("get db returned unexpected error: [%s]", err.Error())
	}
	// this db should be marked as seedable as its new
	if !sq.Seed().Seedable() {
		t.Error("database should be marked as seedable as its new.")
	}
	// check that setting it and then rechecking works
	sq.Seed().Seeded()
	if sq.Seed().Seedable() {
		t.Error("database has been marked as seeded, should fail.")
	}
	// check transaction creation

	if _, err = sq.TX().Get(ctx, sq.DB(), sq.Connector(), sq.Mode()); err != nil {
		t.Errorf("unexpected error getting transaction: [%s]", err.Error())
	}
	// check committing an empty transaction
	if err = sq.TX().Commit(false); err != nil {
		t.Errorf("unexpected error committing transactions [%s]", err.Error())
	}
	// check the date formats are used
	df := sq.Format().Date(dateintervals.Day)
	if df != dateformats.SqliteYMD {
		t.Errorf("date format mismatch - expected [%s] actual [%v]", dateformats.SqliteYMD, df)
	}
	// check it defaults to ym
	df = sq.Format().Date("foobar")
	if df != dateformats.SqliteYM {
		t.Errorf("date format mismatch - expected [%s] actual [%v]", dateformats.SqliteYMD, df)
	}
}
