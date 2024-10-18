package awscosts_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/datastore/awscosts"
)

// TestDatastoreAwsCostsCreateTable creates a new database
// and the table for awscosts and checks the pragma info
// to ensure its been created
func TestDatastoreAwsCostsCreateTable(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "created-table.db")
	var ctx context.Context = context.Background()

	db, err = datastore.New(ctx, dbFile)
	defer db.Close()
	defer os.Remove(dbFile)

	if err != nil {
		t.Errorf("unexpected error creating new database (%s): [%s]", dbFile, err.Error())
	}

	awscosts.Create(ctx, db)

	// check pragma
	pragma := "SELECT name from pragma_table_info('aws_costs')"
	rows, err := db.QueryxContext(ctx, pragma)
	if err != nil {
		t.Errorf("error checking pragma for created table: [%s]", err.Error())
	}

	count := 0
	for rows.Next() {
		count += 1
	}
	if count == 0 {
		t.Errorf("fields for table not found - count: [%d]", count)
	}
}
