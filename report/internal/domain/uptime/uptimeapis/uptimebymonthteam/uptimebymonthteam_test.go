package uptimebymonthteam

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/seed"
	"opg-reports/report/internal/utils/times"
	"opg-reports/report/internal/utils/unmarshal"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jmoiron/sqlx"
)

var sqlParams string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000&_temp_store=memory"

func TestDomainUptimeApiByTeam(t *testing.T) {
	var (
		err      error
		db       *sqlx.DB
		api      humatest.TestAPI
		resp     *httptest.ResponseRecorder
		endpoint string                         = ENDPOINT
		dir      string                         = t.TempDir()
		ctx      context.Context                = t.Context()
		log      *slog.Logger                   = logger.New("error", "text")
		driver   string                         = "sqlite3"
		pth      string                         = filepath.Join(dir, "test-uptime-api.db")
		connStr  string                         = fmt.Sprintf("%s%s", pth, sqlParams)
		apiData  *UptimeByMonthTeamResponseBody = &UptimeByMonthTeamResponseBody{}
		start    string                         = times.AsYMString(times.Add(time.Now(), -8, times.MONTH))
		end      string                         = times.AsYMString(times.Add(time.Now(), -2, times.MONTH))
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

	// register endpoints
	Register(ctx, log, db, api)
	// call the api and get a response, update the endpoint with test values
	endpoint = strings.ReplaceAll(endpoint, "{start_date}", start)
	endpoint = strings.ReplaceAll(endpoint, "{end_date}", end)
	resp = api.GetCtx(ctx, endpoint)
	// unmarshal the data result
	err = unmarshal.FromResponse(resp.Result(), &apiData)
	if err != nil {
		t.Errorf("unexpected unmarshal issue:\n%v", err.Error())
	}
	// test the data...
	if apiData.Count < 1 {
		t.Errorf("expected result count to be higher.")
	}
	if len(apiData.Data) != apiData.Count {
		t.Errorf("data length and count dont match.")
	}

}
