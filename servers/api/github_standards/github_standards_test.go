package github_standards_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
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
	// setup db
	server := api.New(ctx, dbF)
	handler := api.Wrap(server, github_standards.ListHandler)
	// -- setup a mock api thats bound to the correct handler func
	mock := testhelpers.MockServer(handler, "warn")
	defer mock.Close()

	hr, err := httphandler.Get("", "", mock.URL)

	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Debug("api call duration", slog.Float64("seconds", hr.Duration), slog.String("url", mock.URL))

	// -- check values of the response
	_, bytes := convert.Stringify(hr.Response)
	response, _ := convert.Unmarshal[*github_standards.GHSResponse](bytes)

	// -- check the counters match with generated number
	counts := response.Counters
	count := counts.Totals.Count
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", counts)
	}

	// -- call with filters to check status is correct
	list := []string{"?archived=true", "?archived=true&team=foo", "?team=foo"}
	for _, l := range list {
		hr, err := httphandler.Get("", "", mock.URL+l)
		if err != nil {
			slog.Error(err.Error())
			log.Fatal(err.Error())
		}

		slog.Debug("api call duration", slog.Float64("seconds", hr.Duration), slog.String("url", mock.URL+l))
		if hr.StatusCode != http.StatusOK {
			t.Errorf("api call failed")
		}
	}

}
