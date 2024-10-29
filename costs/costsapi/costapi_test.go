package costsapi

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/costs"
	"github.com/ministryofjustice/opg-reports/costs/costsdb"
	"github.com/ministryofjustice/opg-reports/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
)

func testDB(ctx context.Context, file string) (db *sqlx.DB) {
	// make the new db at this location
	db, _, _ = datastore.NewDB(ctx, datastore.Sqlite, file)
	creates := []datastore.CreateStatement{costsdb.CreateCostTable, costsdb.CreateCostTableIndex}
	datastore.Create(ctx, db, creates)
	return
}

func testDBSeed(ctx context.Context, items []*costs.Cost, db *sqlx.DB) *sqlx.DB {
	datastore.InsertMany(ctx, db, costsdb.InsertCosts, items)
	defer db.Close()
	return db
}

// TestCostApiTotalHandler bootstraps a known database
// and then calls the apiTotal func directly to ensure the
// logic is correct and total sum is a match
func TestCostApiTotalHandler(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input *TotalInput = &TotalInput{}
	var dummy = []*costs.Cost{}
	var total float64 = 0.0
	var result *TotalResult

	// generate dummy data set
	dummy = exfaker.Many[costs.Cost](100)
	for _, row := range dummy {
		if row.Service != "Tax" {
			total += row.Value()
		}
	}
	testDBSeed(ctx, dummy, testDB(ctx, dbFile))

	// setup input
	input.Version = "v1"
	input.StartDate = exfaker.TimeStringMin.AddDate(0, 0, -1).Format(exfaker.DateStringFormat)
	input.EndDate = exfaker.TimeStringMax.AddDate(0, 0, 2).Format(exfaker.DateStringFormat)

	result, err = apiTotal(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}
	expected := fmt.Sprintf("%.4f", total)
	actual := fmt.Sprintf("%.4f", result.Body.Result)
	if expected != actual {
		t.Errorf("totals dont match - expected [%s] actual [%s]", expected, actual)
	}
}

// TestCostApiTaxOverviewHandler bootstraps a known database
// and then calls the apiTaxOverview func directly to ensure the
// logic is correct and values for with & without tax match
// the source data
func TestCostApiTaxOverviewHandler(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input *TaxOverviewInput = &TaxOverviewInput{}
	var dummy = []*costs.Cost{}
	var result *TaxOverviewResult

	var month = "2024-01-01"
	var excTax float64 = 0.0
	var incTax float64 = 0.0
	// generate dummy data set - fixed to a single month
	dummy = exfaker.Many[costs.Cost](30)
	for _, row := range dummy {
		incTax += row.Value()
		row.Date = month
		if row.Service != "Tax" {
			excTax += row.Value()
		}
	}

	testDBSeed(ctx, dummy, testDB(ctx, dbFile))

	// setup input
	input.Version = "v1"
	input.Interval = "monthly"
	input.DateFormat = datastore.Sqlite.YearMonthFormat
	input.StartDate = exfaker.TimeStringMin.AddDate(0, 0, -1).Format(exfaker.DateStringFormat)
	input.EndDate = exfaker.TimeStringMax.AddDate(0, 0, 2).Format(exfaker.DateStringFormat)

	result, err = apiTaxOverview(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}

	for _, r := range result.Body.Result {
		var actual = fmt.Sprintf("%.4f", r.Value())
		var expected = ""
		if r.Service == "Including Tax" {
			expected = fmt.Sprintf("%.4f", incTax)
		} else {
			expected = fmt.Sprintf("%.4f", excTax)
		}
		if expected != actual {
			t.Errorf("tax mismatch [%s] - expected [%s] actual [%s]", r.Service, expected, actual)
		}
	}

}

// TestCostApiPerUnitHandler bootstraps a known database
// and then calls the apiPerUnit func directly to ensure the
// logic is correct and values for each unit match
func TestCostApiPerUnitHandler(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input *StandardInput = &StandardInput{}
	var dummy = []*costs.Cost{}
	var result *StandardResult

	var month = "2024-01-01"
	// known units from the faker setup
	var totals map[string]float64 = map[string]float64{
		"unitA": 0.0, "unitB": 0.0, "unitC": 0.0,
	}

	// generate dummy data set - fixed to a single month & ignore tax
	dummy = exfaker.Many[costs.Cost](50)
	for _, row := range dummy {
		row.Date = month
		if row.Service == "Tax" {
			row.Service = "ecs"
		}
		totals[row.Unit] += row.Value()
	}

	testDBSeed(ctx, dummy, testDB(ctx, dbFile))

	// setup input for typical call
	input.Version = "v1"
	input.Interval = "monthly"
	input.DateFormat = datastore.Sqlite.YearMonthFormat
	input.StartDate = exfaker.TimeStringMin.AddDate(0, 0, -1).Format(exfaker.DateStringFormat)
	input.EndDate = exfaker.TimeStringMax.AddDate(0, 0, 2).Format(exfaker.DateStringFormat)

	result, err = apiPerUnit(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}

	for _, r := range result.Body.Result {
		var expected = fmt.Sprintf("%.4f", totals[r.Unit])
		var actual = fmt.Sprintf("%.4f", r.Value())

		if expected != actual {
			t.Errorf("unit [%s] total mismatch - expected [%s] actual [%s]", r.Unit, expected, actual)
		}
	}

	// now filter the results by just one unit
	input.Unit = "unitA"

	result, err = apiPerUnit(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}

	if len(result.Body.Result) != 1 {
		t.Errorf("should only have 1 result, found more.")
	}

	var expected = fmt.Sprintf("%.4f", totals[input.Unit])
	var actual = fmt.Sprintf("%.4f", result.Body.Result[0].Value())
	if expected != actual {
		t.Errorf("unit filtered [%s] result mismatch - expected [%s] actual [%s]", input.Unit, expected, actual)
	}

}

// TestCostApiPerUnitEnvHandler bootstraps a known database
// and then calls the apiPerUnitEnv func directly to ensure the
// logic is correct and values for each unit match
func TestCostApiPerUnitEnvHandler(t *testing.T) {
	var err error
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input *StandardInput = &StandardInput{}
	var dummy = []*costs.Cost{}
	var result *StandardResult

	var month = "2024-01-01"
	// known units from the faker setup
	var totals map[string]map[string]float64 = map[string]map[string]float64{
		"unitA": {
			"production": 0.0, "pre-production": 0.0, "development": 0.0,
		},
		"unitB": {
			"production": 0.0, "pre-production": 0.0, "development": 0.0,
		},
		"unitC": {
			"production": 0.0, "pre-production": 0.0, "development": 0.0,
		},
	}

	// generate dummy data set - fixed to a single month & ignore tax
	dummy = exfaker.Many[costs.Cost](50)
	for _, row := range dummy {
		row.Date = month
		if row.Service == "Tax" {
			row.Service = "ecs"
		}
		totals[row.Unit][row.Environment] += row.Value()
	}

	testDBSeed(ctx, dummy, testDB(ctx, dbFile))

	// setup input for typical call
	input.Version = "v1"
	input.Interval = "monthly"
	input.DateFormat = datastore.Sqlite.YearMonthFormat
	input.StartDate = exfaker.TimeStringMin.AddDate(0, 0, -1).Format(exfaker.DateStringFormat)
	input.EndDate = exfaker.TimeStringMax.AddDate(0, 0, 2).Format(exfaker.DateStringFormat)

	result, err = apiPerUnitEnv(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}

	for _, r := range result.Body.Result {
		var expected = fmt.Sprintf("%.4f", totals[r.Unit][r.Environment])
		var actual = fmt.Sprintf("%.4f", r.Value())

		if expected != actual {
			t.Errorf("unit [%s-%s] total mismatch - expected [%s] actual [%s]", r.Unit, r.Environment, expected, actual)
		}
	}

	// now filter the results by just one unit
	input.Unit = "unitA"

	result, err = apiPerUnitEnv(ctx, input)
	if err != nil {
		t.Errorf("unexpected error getting total: [%s]", err.Error())
	}

	if len(result.Body.Result) <= 0 {
		t.Errorf("should only have at least result")
	}

	for _, r := range result.Body.Result {
		var expected = fmt.Sprintf("%.4f", totals[input.Unit][r.Environment])
		var actual = fmt.Sprintf("%.4f", r.Value())
		if expected != actual {
			t.Errorf("unit filtered [%s-%s] result mismatch - expected [%s] actual [%s]", input.Unit, r.Environment, expected, actual)
		}
	}

}

// TestCostApiRegister checks the register function maps correctly
func TestCostApiRegister(t *testing.T) {

	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var dummy = exfaker.Many[costs.Cost](50)
	var urls = []string{
		"/v1/costs/aws/total/2024-01-01/2024-01-01",
		"/v1/costs/aws/tax-overview/2024-01-01/2024-01-01/month",
		"/v1/costs/aws/tax-overview/2024-01-01/2024-01-01/day",
		"/v1/costs/aws/unit/2024-01-01/2024-01-01/month",
		"/v1/costs/aws/unit/2024-01-01/2024-01-01/month?unit=unitA",
		"/v1/costs/aws/unit/2024-01-01/2024-01-01/day",
		"/v1/costs/aws/unit-environment/2024-01-01/2024-01-01/month",
		"/v1/costs/aws/unit-environment/2024-01-01/2024-01-01/day",
		"/v1/costs/aws/unit-environment/2024-01-01/2024-01-01/month?unit=A",
		"/v1/costs/aws/detailed/2024-01-01/2024-01-01/month",
		"/v1/costs/aws/detailed/2024-01-01/2024-01-01/day",
		"/v1/costs/aws/detailed/2024-01-01/2024-01-01/month?unit=A",
	}
	var middleware = func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, Segment, dbFile)
		next(ctx)
	}

	// t.Log()
	testDBSeed(ctx, dummy, testDB(ctx, dbFile))

	// register the routes
	_, api := humatest.New(t, huma.DefaultConfig("Reporting API", "test"))
	api.UseMiddleware(middleware)
	Register(api)

	for _, uri := range urls {
		resp := api.Get(uri)
		if resp.Code != http.StatusOK {
			t.Errorf("endpoint [%s] failed with code [%v]", uri, resp.Code)
		}
	}

}
