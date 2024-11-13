package sqlite_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors/sqlite"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

// interface testing
var (
	_ dbs.Connector     = &sqlite.Sqlite{}
	_ dbs.Formattable   = &sqlite.Sqlite{}
	_ dbs.Seeder        = &sqlite.Sqlite{}
	_ dbs.Transactional = &sqlite.Sqlite{}
	_ dbs.Adaptor       = &sqlite.Sqlite{}
)

func TestSqliteSetup(t *testing.T) {
	var (
		err        error
		sq         *sqlite.Sqlite
		connect    dbs.ConnectionString
		conn       string
		ctx        context.Context = context.Background()
		driver     string          = "sqlite3"
		paramChunk string          = "_journal=WAL"
		dir        string          = t.TempDir()
		path       string          = filepath.Join(dir, "test.db")
	)

	sq, err = sqlite.New(path)
	if err != nil {
		t.Errorf("error calling new: [%s]", err.Error())
	}
	// test it created a db file
	if !fileutils.Exists(path) {
		t.Errorf("database file was not created")
	}
	// test driver name - this is fixed on creation, so should always match
	if driver != string(sq.GetDriverName()) {
		t.Errorf("drivers did not match - expected [%s] actual [%v]", driver, sq.GetDriverName())
	}
	// test path name is as we set it
	if path != string(sq.GetPath()) {
		t.Errorf("paths did not match - expected [%s] actual [%v]", path, sq.GetPath())
	}
	// test params have the wal in them
	if !strings.Contains(string(sq.GetParams()), paramChunk) {
		t.Errorf("params did not match - expected to contains [%s] actual [%v]", paramChunk, sq.GetParams())
	}
	// test the connection string has both path and param chunk
	connect = sq.GetConnectionString(sq.GetPath(), sq.GetParams())
	conn = string(connect)
	if !strings.Contains(conn, path) {
		t.Errorf("connection did not match - expected to contains [%s] actual [%v]", path, conn)
	}
	if !strings.Contains(conn, paramChunk) {
		t.Errorf("connection did not match - expected to contains [%s] actual [%v]", paramChunk, conn)
	}
	// check the db fetch worked
	if _, err = sq.GetDB(ctx, sq.GetDriverName(), connect); err != nil {
		t.Errorf("get db returned unexpected error: [%s]", err.Error())
	}
	// this db should be marked as seedable as its new
	if !sq.Seedable() {
		t.Error("database should be marked as seedable as its new.")
	}
	// check that setting it and then rechecking works
	sq.Seeded()
	if sq.Seedable() {
		t.Error("database has been marked as seeded, should fail.")
	}
	// check transaction creation
	if _, err = sq.GetTransaction(ctx, sq.MustGetDB(ctx, sq.GetDriverName(), connect), true); err != nil {
		t.Errorf("unexpected error getting transaction: [%s]", err.Error())
	}
	// check committing an empty transaction
	if err = sq.CommitTransaction(sq.MustGetTransaction(ctx, sq.MustGetDB(ctx, sq.GetDriverName(), connect), true), false); err != nil {
		t.Errorf("unexpected error committing transactions [%s]", err.Error())
	}
	// check the date formats are used
	df := sq.DateFormat(dateintervals.Day)
	if df != dateformats.SqliteYMD {
		t.Errorf("date format mismatch - expected [%s] actual [%v]", dateformats.SqliteYMD, df)
	}
	// check it defaults to ym
	df = sq.DateFormat("foobar")
	if df != dateformats.SqliteYM {
		t.Errorf("date format mismatch - expected [%s] actual [%v]", dateformats.SqliteYMD, df)
	}
}
