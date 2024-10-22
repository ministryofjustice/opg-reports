package awscosts_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/datastore/awscosts"
	"github.com/ministryofjustice/opg-reports/fakes"
)

// TestDatastoreAwsCostsQueriesTotalWithinDateRange creates and then
// inserts a series of sample data and then checks the total query
// generated matches the total from the sample data
// Tests Single function
func TestDatastoreAwsCostsQueriesTotalWithinDateRange(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "total.db")
	var ctx context.Context = context.Background()
	var insertCount int = 10
	var inserts []*awscosts.Cost = awscosts.Fakes(insertCount)
	var expectedTotal float64 = 0.0
	var actualTotal float64 = 0.0

	db, _, err = datastore.New(ctx, datastore.Sqlite, dbFile)
	defer db.Close()
	defer os.Remove(dbFile)

	if err != nil {
		t.Errorf("unexpected error creating new database (%s): [%s]", dbFile, err.Error())
	}

	awscosts.Create(ctx, db)

	// -- insert the faked items
	_, err = awscosts.InsertMany(ctx, db, inserts)
	if err != nil {
		t.Errorf("failed to insert multiple records:\n [%s]", err.Error())
	}

	// work out the expected total
	for _, faked := range inserts {
		expectedTotal += faked.Value()
	}

	result, err := awscosts.GetOne(ctx, db, awscosts.TotalInDateRange, fakes.MinDate, fakes.MaxDate)
	if err != nil {
		t.Errorf("error from getting total: [%s]", err.Error())
	}
	actualTotal = result.(float64)

	if fmt.Sprintf("%.4f", actualTotal) != fmt.Sprintf("%.4f", expectedTotal) {
		t.Errorf("total does not match expected - expected [%v] actual [%v]", expectedTotal, actualTotal)
	}

}

// TestDatastoreAwsCostsQueriesTotalsWithAndWithoutTax creates and seeds
// a dummy database then runs query for TotalsWithAndWithoutTax to ensure
// to totals map to the sample data
// Tests Many function
func TestDatastoreAwsCostsQueriesTotalsWithAndWithoutTax(t *testing.T) {

	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	// var dir string = "./"
	var dbFile string = filepath.Join(dir, "with-without-tax.db")
	var ctx context.Context = context.Background()

	db, _, err = datastore.New(ctx, datastore.Sqlite, dbFile)
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
	_, err = awscosts.InsertMany(ctx, db, inserts)
	if err != nil {
		t.Errorf("unexpected error inserting data: [%s]", err.Error())
	}

	// -- run the query for a month
	params := &awscosts.NamedParameters{
		StartDate:  "2024-01-01",
		EndDate:    "2024-04-01",
		DateFormat: datastore.Sqlite.YearMonthFormat,
	}
	results, err := awscosts.GetMany(ctx, db, awscosts.TotalsWithAndWithoutTax, params)
	if err != nil {
		t.Errorf("unxpected error on query: [%s]", err.Error())
	}

	// -- there should be 6 rows (1 with, 1 without x 3 months)
	expectedCount := 6
	many := results
	if len(many) != expectedCount {
		t.Errorf("expected [%d] rows, actual [%v]", expectedCount, len(many))
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
	for _, res := range many {
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

// TestDatastoreAwsCostsNeeds checks adds named parameters to
// a string and checks the Needs returns the correct values.
// Needs is used in validating named parameters being set for
// certain queries (ManyStatement)
func TestDatastoreAwsCostsNeeds(t *testing.T) {
	var testString awscosts.ManyStatement = ""
	var found []string = []string{}

	// this should return min & max
	testString = `SELECT count(*) FROM table WHERE id > :min AND id < :max`
	found = awscosts.Needs(testString)
	if !slices.Contains(found, "min") || !slices.Contains(found, "max") {
		t.Errorf("needs did not find all fields: [%s]", strings.Join(found, ","))
	}
	if len(found) != 2 {
		t.Errorf("incorrect number of fields found: [%s]", strings.Join(found, ","))
	}
	// testing multiline string
	testString = `
	SELECT
		count(*)
	FROM table
	WHERE
		id > :min
		AND id < :max
		AND date_end is TRUE
		AND start > :start_dateWith-oddName
	`
	found = awscosts.Needs(testString)
	if !slices.Contains(found, "min") || !slices.Contains(found, "max") || !slices.Contains(found, "start_dateWith-oddName") {
		t.Errorf("needs did not find all fields: [%s]", strings.Join(found, ","))
	}
	if len(found) != 3 {
		t.Errorf("incorrect number of fields found: [%s]", strings.Join(found, ","))
	}

}

// TestDatastoreAwsCostsValidateParameters checks that valid
// param and needs as well as invalid versions are triggered
// correctly
func TestDatastoreAwsCostsValidateParameters(t *testing.T) {
	var err error
	var needs []string = []string{}
	var params *awscosts.NamedParameters = &awscosts.NamedParameters{}

	// -- test it works and ignores extra fields
	needs = []string{"start_date"}
	params = &awscosts.NamedParameters{StartDate: "test", EndDate: "test"}
	err = awscosts.ValidateParameters(params, needs)
	if err != nil {
		t.Errorf("param should be valid: [%s]", err.Error())
	}

	// -- test a failing one
	needs = []string{"end_date"}
	params = &awscosts.NamedParameters{StartDate: "test"}
	err = awscosts.ValidateParameters(params, needs)
	if err == nil {
		t.Errorf("param should throw error, but didnt")
	}

}
