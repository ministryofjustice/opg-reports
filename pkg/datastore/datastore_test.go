package datastore_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
)

// TestDatastoreNewCreatesDbFile checks that datastore.New
// successfully creates an empty database if the file
// passed does not curretly exist
func TestDatastoreNewCreatesDbFile(t *testing.T) {
	var err error
	var db *sqlx.DB
	var isNew bool = true

	dir := t.TempDir()
	file := filepath.Join(dir, "file-does-not-exist.db")
	defer os.Remove(file)

	db, isNew, err = datastore.NewDB(context.Background(), datastore.Sqlite, file)
	defer db.Close()
	if err != nil {
		t.Errorf("error from datastore.New: %s", err.Error())
	}

	if !isNew {
		t.Errorf("new database should have returned as being new")
	}
	// fail if there is an error stating the file
	if _, err = os.Stat(file); err != nil {
		t.Errorf("datastore.New did not create file (%s): [%s]", file, err.Error())
	}
}

// TestDatastoreNewPing checks that the db returned from
// datastore.New pings successfully
func TestDatastoreNewPing(t *testing.T) {
	var err error
	var ctx context.Context = context.Background()

	dir := t.TempDir()
	file := filepath.Join(dir, "ping.db")
	defer os.Remove(file)

	db, _, err := datastore.NewDB(ctx, datastore.Sqlite, file)
	defer db.Close()

	if err != nil {
		t.Errorf("error from datastore.New: %s", err.Error())
	}

	if err = db.PingContext(ctx); err != nil {
		t.Errorf("db.PingContext throw an error: [%s]", err.Error())
	}

}