package awscosts_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/awscosts"
	"github.com/ministryofjustice/opg-reports/datastore"
)

// TestDatastoreAwsCostsInsertItems creates a test database
// and then inserts a series of faked cost records and then
// confirms the table count to ensure all were created
func TestDatastoreAwsCostsInsertItems(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "create-and-insert.db")
	var ctx context.Context = context.Background()
	var insertCount int = 5
	var inserts []*awscosts.Cost = awscosts.Fakes(insertCount)

	db, _, err = datastore.New(ctx, datastore.Sqlite, dbFile)
	defer db.Close()
	defer os.Remove(dbFile)

	if err != nil {
		t.Errorf("unexpected error creating new database (%s): [%s]", dbFile, err.Error())
	}
	awscosts.Create(ctx, db)

	for _, item := range inserts {
		_, err := awscosts.InsertOne(ctx, db, item)
		if err != nil {
			t.Errorf("failed to insert awscost item: [%s]", err.Error())
		}
	}

	countSql := `SELECT count(*) from aws_costs`
	found := 0
	err = db.Get(&found, countSql)
	if err != nil {
		t.Errorf("failed to count entries in aws_costs: [%s]", err.Error())
	}
	if found != insertCount {
		t.Errorf("incorrect number of entries found - expected [%d] actual [%v]", insertCount, found)
	}

}

// TestDatastoreAwsCostsInsertAll creates a test database
// and then inserts a series of faked cost records at once
// using InsertAll
func TestDatastoreAwsCostsInsertAll(t *testing.T) {
	var err error
	var db *sqlx.DB
	// var dir string = "./"
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "create-and-insert-all.db")
	var ctx context.Context = context.Background()
	// -- 1.5 million rows for perf testing
	// var insertCount int = 1500000
	// -- 15k for faster tests
	var insertCount int = 15000
	var inserts []*awscosts.Cost = awscosts.Fakes(insertCount)
	var ids []int = []int{}

	db, _, err = datastore.New(ctx, datastore.Sqlite, dbFile)
	defer db.Close()
	defer os.Remove(dbFile)

	if err != nil {
		t.Errorf("unexpected error creating new database (%s): [%s]", dbFile, err.Error())
	}
	awscosts.Create(ctx, db)

	ids, err = awscosts.InsertMany(ctx, db, inserts)
	if err != nil {
		t.Errorf("failed to insert multiple records:\n [%s]", err.Error())
	}

	countSql := `SELECT count(*) from aws_costs`
	found := 0
	err = db.Get(&found, countSql)
	if err != nil {
		t.Errorf("failed to count entries in aws_costs: [%s]", err.Error())
	}
	if len(ids) != insertCount {
		t.Errorf("incorrect number of id's returned - expected [%d] actual [%v]", insertCount, len(ids))
	}

	if found != insertCount {
		t.Errorf("incorrect number of entries found - expected [%d] actual [%v]", insertCount, found)
	}

}
