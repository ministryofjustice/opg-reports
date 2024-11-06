package standardsapi

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/sources/standards"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsdb"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsio"
)

func testDB(ctx context.Context, file string) (db *sqlx.DB) {
	// make the new db at this location
	db, _, _ = datastore.NewDB(ctx, datastore.Sqlite, file)
	creates := []datastore.CreateStatement{
		standardsdb.CreateStandardsTable,
		standardsdb.CreateStandardsIndexIsArchived,
		standardsdb.CreateStandardsIndexIsArchivedTeams,
		standardsdb.CreateStandardsIndexTeams,
		standardsdb.CreateStandardsIndexBaseline,
		standardsdb.CreateStandardsIndexExtended,
	}

	datastore.Create(ctx, db, creates)
	return
}

func testDBSeed(ctx context.Context, items []*standards.Standard, db *sqlx.DB) *sqlx.DB {
	datastore.InsertMany(ctx, db, standardsdb.InsertStandard, items)
	defer db.Close()
	return db
}

// TestStandardsApiHandlerArchive check the api returns correct number
// of archived results
func TestStandardsApiHandlerArchive(t *testing.T) {
	var err error
	var out *standardsio.StandardsOutput

	var n = 100
	var expectedArchivedCount = 0
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input = &standardsio.StandardsInput{Version: "v1", Archived: true}
	var fakes = []*standards.Standard{}

	fakes = exfaker.Many[*standards.Standard](n)
	testDBSeed(ctx, fakes, testDB(ctx, dbFile))

	for _, f := range fakes {
		if f.Archived() {
			expectedArchivedCount += 1
		}
	}

	out, err = apiHandler(ctx, input)
	if err != nil {
		t.Errorf("error in handler: [%s]", err.Error())
	}

	if out.Body.Counters.Total != n {
		t.Errorf("db content mismatch")
	}

	if len(out.Body.Result) != expectedArchivedCount {
		t.Errorf("api handler returned incorrect number of archived entries")
	}

}

// TestStandardsApiHandlerTeamFilter check the api returns correct number
// of results for the team
func TestStandardsApiHandlerTeamFilter(t *testing.T) {
	var err error
	var out *standardsio.StandardsOutput

	var n = 100
	var team = "#unitA#" // one of the default teams used by faker
	var expected = 0
	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var input = &standardsio.StandardsInput{Version: "v1", Archived: false, Unit: "unitA"}
	var fakes = []*standards.Standard{}

	fakes = exfaker.Many[*standards.Standard](n)
	testDBSeed(ctx, fakes, testDB(ctx, dbFile))

	for _, f := range fakes {
		if f.Teams == team && !f.Archived() {
			expected += 1
		}
	}
	// make sure to call the resolve to create .Teams from Unit
	input.Resolve(nil)
	out, err = apiHandler(ctx, input)
	if err != nil {
		t.Errorf("error in handler: [%s]", err.Error())
	}

	if out.Body.Counters.Total != n {
		t.Errorf("db content mismatch")
	}
	actual := len(out.Body.Result)
	if actual != expected {
		t.Errorf("api handler returned incorrect number of filtered, unarchived entries - expected [%d] actual [%v]", expected, actual)
	}

}

// TestStandardsApiRegister checks that the url returns a correct status code
func TestStandardsApiRegister(t *testing.T) {

	var dir = t.TempDir()
	var dbFile = filepath.Join(dir, "test.db")
	var ctx = context.WithValue(context.Background(), Segment, dbFile)
	var dummy = exfaker.Many[*standards.Standard](50)
	var urls = []string{
		"/v1/standards/github/false",
		"/v1/standards/github/false",
		"/v1/standards/github/false?unit=unitA",
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
