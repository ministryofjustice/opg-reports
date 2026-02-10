package codeownerforteam

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/seed"
	"opg-reports/report/internal/utils/unmarshal"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jmoiron/sqlx"
)

func TestDomainCodeownerApiForTeam(t *testing.T) {
	var (
		err     error
		db      *sqlx.DB
		api     humatest.TestAPI
		resp    *httptest.ResponseRecorder
		dir     string                        = t.TempDir()
		ctx     context.Context               = t.Context()
		log     *slog.Logger                  = logger.New("error", "text")
		driver  string                        = "sqlite3"
		connStr string                        = fmt.Sprintf("%s/%s", dir, "test-codeownerforteam-api.db")
		apiData *CodeownerForTeamResponseBody = &CodeownerForTeamResponseBody{}
		ep      string                        = ENDPOINT
	)
	// setup the test huma instance
	_, api = humatest.New(t)
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		t.Errorf("unexpected connection issue:\n%v", err.Error())
	}
	defer db.Close()
	// seed & migrate db
	err = seed.SeedDB(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected seed issue:\n%v", err.Error())
	}
	// add the filter
	ep = strings.ReplaceAll(ep, "{team}", "TEAM-B")
	// register endpoints
	Register(ctx, log, db, api)
	// call the api and get a response
	resp = api.GetCtx(ctx, ep)
	// unmarshal the data result
	err = unmarshal.FromResponse(resp.Result(), &apiData)
	if err != nil {
		t.Errorf("unexpected unmarshal issue:\n%v", err.Error())
	}
	// test the data...
	if apiData.Count < 2 {
		t.Errorf("expected result count to be higher.")
	}
	if len(apiData.Data) != apiData.Count {
		t.Errorf("data length and count dont match.")
	}
}
