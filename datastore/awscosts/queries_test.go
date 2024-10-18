package awscosts_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/datastore/awscosts"
)

// TestDatastoreAwsCostsQueriesTotalsWithAndWithoutTax creates and seeds
// a dummy database then runs query for TotalsWithAndWithoutTax to ensure
// to totals map to the sample data
func TestDatastoreAwsCostsQueriesTotalsWithAndWithoutTax(t *testing.T) {

	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	// var dir string = "./"
	var dbFile string = filepath.Join(dir, "with-without-tax.db")
	var ctx context.Context = context.Background()

	db, err = datastore.New(ctx, dbFile)
	defer db.Close()
	defer os.Remove(dbFile)

	if err != nil {
		t.Errorf("unexpected error creating new database (%s): [%s]", dbFile, err.Error())
	}

	awscosts.Create(ctx, db)

	// -- sample data
	inserts := []*awscosts.Cost{
		{
			Ts:           time.Now().UTC().Format(time.RFC3339),
			Organisation: "test",
			AccountID:    "01",
			AccountName:  "One",
			Unit:         "teamOne",
			Label:        "team one prod",
			Environment:  "production",
			Region:       "us-east-1",
			Service:      "EC2",
			Date:         "2024-01-01",
			Cost:         "10.01",
		},
		{
			Ts:           time.Now().UTC().Format(time.RFC3339),
			Organisation: "test",
			AccountID:    "01",
			AccountName:  "One",
			Unit:         "teamOne",
			Label:        "team one prod",
			Environment:  "production",
			Region:       "us-east-1",
			Service:      "EC2",
			Date:         "2024-02-01",
			Cost:         "12.34",
		},
		{
			Ts:           time.Now().UTC().Format(time.RFC3339),
			Organisation: "test",
			AccountID:    "01",
			AccountName:  "One",
			Unit:         "teamOne",
			Label:        "team one prod",
			Environment:  "production",
			Region:       "us-east-1",
			Service:      "Tax",
			Date:         "2024-02-01",
			Cost:         "1.234",
		},
		{
			Ts:           time.Now().UTC().Format(time.RFC3339),
			Organisation: "test",
			AccountID:    "01",
			AccountName:  "One",
			Unit:         "teamOne",
			Label:        "team one prod",
			Environment:  "production",
			Region:       "us-east-1",
			Service:      "EC2",
			Date:         "2024-03-01",
			Cost:         "55.07",
		},
		{
			Ts:           time.Now().UTC().Format(time.RFC3339),
			Organisation: "test",
			AccountID:    "01",
			AccountName:  "One",
			Unit:         "teamOne",
			Label:        "team one prod",
			Environment:  "production",
			Region:       "us-east-1",
			Service:      "Tax",
			Date:         "2024-03-01",
			Cost:         "7.15",
		},
	}
	_, err = awscosts.InsertAll(ctx, db, inserts)
	if err != nil {
		t.Errorf("unexpected error inserting data: [%s]", err.Error())
	}

	// -- run the query for a month
	params := &awscosts.Parameters{
		StartDate:  "2024-01-01",
		EndDate:    "2024-04-01",
		DateFormat: datastore.Sqlite.YearMonthFormat,
	}
	results, err := awscosts.Query(ctx, db, awscosts.TotalsWithAndWithoutTax, params)
	if err != nil {
		t.Errorf("unxpected error on query: [%s]", err.Error())
	}

	// -- there should be 6 rows (1 with, 1 without x 3 months)
	expectedCount := 6
	if len(results) != expectedCount {
		t.Errorf("expected [%d] rows, actual [%v]", expectedCount, len(results))
	}

	// -- totals for each month without tax
	totalsNoTax := map[string]string{
		"2024-01": "10.01",
		"2024-02": "12.34",
		"2024-03": "55.07",
	}
	totalsWithTax := map[string]string{
		"2024-01": "10.01",
		"2024-02": "13.574",
		"2024-03": "62.22",
	}

	matched := 0
	for _, res := range results {
		comp := totalsNoTax
		key := res.Date
		if res.Service == "Including Tax" {
			comp = totalsWithTax
		}
		if comp[key] != res.Cost {
			t.Errorf("error in tax calc -[%s] month [%s] expected [%s] actual [%s]", res.Service, key, comp[key], res.Cost)
			fmt.Printf("%#v", res)
		} else {
			matched += 1
		}
	}

	if matched != expectedCount {
		t.Errorf("one tax details failed in data")
	}

}
