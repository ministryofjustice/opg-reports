package github_standards_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema string = "../../../datastore/github_standards/github_standards.sql"

// test that creating and then callign an endpoint for github standards returns all
// the data we would expect
func TestServersApiGithubStandardsArchivedApiCallAndParse(t *testing.T) {
	logger.LogSetup()
	ctx := context.TODO()
	N := 500
	dir := t.TempDir()

	dbF := filepath.Join(dir, "ghs.db")
	schemaF := filepath.Join(dir, "ghs.sql")
	dataF := filepath.Join(dir, "dummy.json")

	testhelpers.CopyFile(realSchema, schemaF)
	tick := testhelpers.T()
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "github_standards", N)
	tick.Stop()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	slog.Debug("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := ghs.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}
	// -- set db
	github_standards.SetDBPath(dbF)
	github_standards.SetCtx(ctx)
	// -- setup a mock api thats bound to the correct handler func
	mock := mockApi()
	defer mock.Close()
	u, err := url.Parse(mock.URL)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	// -- call the api - time its duration
	tick = testhelpers.T()
	hr, err := getter.GetUrl(u)
	tick.Stop()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Debug("api call duration", slog.String("seconds", tick.Seconds()), slog.String("u", u.String()))

	// -- check values of the response
	_, bytes := convert.Stringify(hr)
	response, _ := convert.Unmarshal[*resp.Response](bytes)

	// -- check the counters match with generated number
	counts := response.Metadata["counters"].(map[string]interface{})
	all := counts["totals"].(map[string]interface{})
	count := int(all["count"].(float64))
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", all)
	}

	// -- call with filters to check status is correct
	list := []string{"?archived=true", "?archived=true&team=foo", "?team=foo"}
	for _, l := range list {
		tick = testhelpers.T()
		call := u.String() + l
		ur, err := url.Parse(call)
		if err != nil {
			slog.Error(err.Error())
			log.Fatal(err.Error())
		}
		hr, err = getter.GetUrl(ur)
		if err != nil {
			slog.Error(err.Error())
			log.Fatal(err.Error())
		}
		tick.Stop()

		slog.Debug("api call duration", slog.String("seconds", tick.Seconds()), slog.String("url", ur.String()))
		if hr.StatusCode != http.StatusOK {
			t.Errorf("api call failed")
		}
	}

}

func mockApi() *httptest.Server {
	return testhelpers.MockServer(github_standards.ListHandler, "warn")
}
