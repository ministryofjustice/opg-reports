package uptimeapi

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimedb"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeio"
)

func testDB(ctx context.Context, file string) (db *sqlx.DB) {
	// make the new db at this location
	db, _, _ = datastore.NewDB(ctx, datastore.Sqlite, file)
	creates := []datastore.CreateStatement{
		uptimedb.CreateUptimeTable,
		uptimedb.CreateUptimeTableDateIndex,
		uptimedb.CreateUptimeTableUnitDateIndex,
	}

	datastore.Create(ctx, db, creates)
	return
}

func testDBSeed(ctx context.Context, items []*uptime.Uptime, db *sqlx.DB) *sqlx.DB {
	datastore.InsertMany(ctx, db, uptimedb.InsertUptime, items)
	defer db.Close()
	return db
}

// TestUptimeApiHandlerOverall check the api returns correct number
// of results for the overall endpoint
// the api should be grouping the data by the yyyy-mm value
// based on the input, so find the number of groups within
// the fake data and check that matches the handler results
func TestUptimeApiHandlerOverall(t *testing.T) {
	var err error
	var out *uptimeio.UptimeOutput
	var start = exfaker.TimeStringMin.Format(consts.DateFormatYearMonthDay)
	var end = exfaker.TimeStringMax.Format(consts.DateFormatYearMonthDay)
	var n = 100
	var expectedCounters = map[string]int{}
	var expectedCount = 0
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input = &uptimeio.UptimeInput{Version: "v1", StartDate: start, EndDate: end, Interval: "month"}
	var fakes = []*uptime.Uptime{}

	fakes = exfaker.Many[*uptime.Uptime](n)
	testDBSeed(ctx, fakes, testDB(ctx, dbFile))

	for _, f := range fakes {
		var group = convert.DateReformat(f.Date, consts.DateFormatYearMonth)
		if _, ok := expectedCounters[group]; !ok {
			expectedCounters[group] = 0
		}
		expectedCounters[group] += 1
	}
	expectedCount = len(expectedCounters)

	input.Resolve(nil)
	out, err = apiOverallHandler(ctx, input)
	if err != nil {
		t.Errorf("error in handler: [%s]", err.Error())
	}

	if len(out.Body.Result) != expectedCount {
		t.Errorf("api handler returned incorrect number of entries - expected [%d] actual [%v]", expectedCount, len(out.Body.Result))
	}

}

// TestUptimeApiHandlerUnit check the api returns correct number
// of results for the unit endpoint
// the api should be grouping the data by the yyyy-mm & unit value
// based on the input, so find the number of groups within
// the fake data and check that matches the handler results
func TestUptimeApiHandlerUnit(t *testing.T) {
	var err error
	var out *uptimeio.UptimeOutput
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var start = exfaker.TimeStringMin.Format(consts.DateFormatYearMonthDay)
	var end = exfaker.TimeStringMax.Format(consts.DateFormatYearMonthDay)
	var n = 100
	var expectedCounters = map[string]int{}
	var expectedCount = 0
	var input = &uptimeio.UptimeInput{Version: "v1", StartDate: start, EndDate: end, Interval: "month"}
	var fakes = []*uptime.Uptime{}

	fakes = exfaker.Many[*uptime.Uptime](n)
	testDBSeed(ctx, fakes, testDB(ctx, dbFile))

	for _, f := range fakes {
		var group = fmt.Sprintf("%s-%s", convert.DateReformat(f.Date, consts.DateFormatYearMonth), f.Unit)
		if _, ok := expectedCounters[group]; !ok {
			expectedCounters[group] = 0
		}
		expectedCounters[group] += 1
	}
	expectedCount = len(expectedCounters)

	input.Resolve(nil)
	out, err = apiUnitHandler(ctx, input)
	if err != nil {
		t.Errorf("error in handler: [%s]", err.Error())
	}

	if len(out.Body.Result) != expectedCount {
		t.Errorf("api handler returned incorrect number of entries - expected [%d] actual [%v]", expectedCount, len(out.Body.Result))
	}

}

// TestUptimeApiHandlerUnitFilter check the api returns correct number
// of results for the unit endpoint when filtering by unit
// the api should be grouping the data by the yyyy-mm & unit value
// based on the input, so find the number of groups within
// the fake data for this unit only and check that matches
// the handler results
func TestUptimeApiHandlerUnitFilter(t *testing.T) {
	var err error
	var out *uptimeio.UptimeOutput
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var start = exfaker.TimeStringMin.Format(consts.DateFormatYearMonthDay)
	var end = exfaker.TimeStringMax.Format(consts.DateFormatYearMonthDay)
	var n = 100
	var expectedCounters = map[string]int{}
	var expectedCount = 0
	var input = &uptimeio.UptimeInput{Version: "v1", StartDate: start, EndDate: end, Interval: "month", Unit: "unitA"}
	var fakes = []*uptime.Uptime{}

	fakes = exfaker.Many[*uptime.Uptime](n)
	testDBSeed(ctx, fakes, testDB(ctx, dbFile))

	for _, f := range fakes {
		if f.Unit == "unitA" {
			var group = fmt.Sprintf("%s-%s", convert.DateReformat(f.Date, consts.DateFormatYearMonth), f.Unit)
			if _, ok := expectedCounters[group]; !ok {
				expectedCounters[group] = 0
			}
			expectedCounters[group] += 1
		}
	}
	expectedCount = len(expectedCounters)

	input.Resolve(nil)
	out, err = apiUnitHandler(ctx, input)
	if err != nil {
		t.Errorf("error in handler: [%s]", err.Error())
	}

	if len(out.Body.Result) != expectedCount {
		t.Errorf("api handler returned incorrect number of entries - expected [%d] actual [%v]", expectedCount, len(out.Body.Result))
	}

}

func TestUptimeApiRegister(t *testing.T) {

	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var dummy = exfaker.Many[*uptime.Uptime](1440)
	var start = exfaker.TimeStringMin.Format(consts.DateFormatYearMonthDay)
	var end = exfaker.TimeStringMax.Format(consts.DateFormatYearMonthDay)
	var urls = []string{
		"/v1/uptime/aws/overall/" + start + "/" + end + "/day",
		"/v1/uptime/aws/overall/" + start + "/" + end + "/month",
		"/v1/uptime/aws/unit/" + start + "/" + end + "/day",
		"/v1/uptime/aws/unit/" + start + "/" + end + "/month",
		"/v1/uptime/aws/unit/" + start + "/" + end + "/month?unit=unitA",
	}
	var middleware = func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, Segment, dbFile)
		next(ctx)
	}

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
