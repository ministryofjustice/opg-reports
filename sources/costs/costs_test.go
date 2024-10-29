package costs_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsdb"
)

// TestCostsModelValue checks the .Value func
// on the Cost struct converts from strings
// to floats correctly
func TestCostsModelValue(t *testing.T) {
	var cost = &costs.Cost{}
	var values = map[string]float64{
		"10.0101":  10.0101,
		"100.0":    100.0,
		"-4.56134": -4.56134,
	}

	for str, expected := range values {
		cost.Cost = str
		actual := cost.Value()
		if expected != actual {
			t.Errorf("float conversion error - expected [%f] actual [%v]", expected, actual)
		}
	}

}

// TestCostsSetup runs Setup and then checks the
// file exists and the location and that it has
// records in the table
func TestCostsSetup(t *testing.T) {
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.Background()

	costs.Setup(ctx, dbFile)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Errorf("database file was not created in expected locations")
	}
	db, _, err := datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("error connecting to db [%s]", err.Error())
	}

	count, err := datastore.Get[int](ctx, db, costsdb.RowCount)
	if err != nil {
		t.Errorf("error counting db rows: [%s]", err.Error())
	}
	if count != costs.RecordsToSeed {
		t.Errorf("incorrect number of rows - expected [%d] actual [%v]", costs.RecordsToSeed, count)
	}
}
