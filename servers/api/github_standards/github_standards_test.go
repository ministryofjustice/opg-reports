package github_standards_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/seeder"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

// seed the database and then run filters to test and view performance
func TestServersApiGithubStandardsArchivedPerfDBOnly(t *testing.T) {
	logger.LogSetup()
	ctx := context.Background()
	N := 50000

	dir := t.TempDir()
	dbName := filepath.Join(dir, "github_standards.db")
	seed := testSeed(dbName, N)

	s := time.Now().UTC()
	db, err := seeder.New(seed)
	defer db.Close()

	e := time.Now().UTC()
	dur := e.Sub(s)

	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}
	slog.Warn("seed duration", slog.Float64("seconds", dur.Seconds()))
	q := ghs.New(db)

	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	s = time.Now().UTC()
	res, _ := q.FilterByIsArchived(ctx, 1)
	e = time.Now().UTC()
	dur = e.Sub(s)
	slog.Warn("archived filter duration",
		slog.Float64("seconds", dur.Seconds()),
		slog.Int("records", len(res)))

	s = time.Now().UTC()
	team := "%#" + "foo" + "#%"
	res, _ = q.FilterByIsArchivedAndTeam(ctx, ghs.FilterByIsArchivedAndTeamParams{IsArchived: 1, Teams: team})
	e = time.Now().UTC()
	dur = e.Sub(s)
	slog.Warn("archived team filter duration",
		slog.Float64("seconds", dur.Seconds()),
		slog.Int("records", len(res)))

}

func TestServersApiGithubStandardsArchivedPerfApiCallAndParse(t *testing.T) {
	logger.LogSetup()
	slog.Warn("start")
	ctx := context.TODO()
	N := 50000
	dir := t.TempDir()
	dbName := filepath.Join(dir, "github_standards.db")
	seed := testSeed(dbName, N)

	s := time.Now().UTC()
	db, err := seeder.New(seed)
	defer db.Close()
	e := time.Now().UTC()
	dur := e.Sub(s)

	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}
	slog.Warn("seed duration", slog.Float64("seconds", dur.Seconds()))

	q := ghs.New(db)
	defer q.Close()

	// slog.Warn("counting")
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}
	// slog.Warn("mocking api")
	mock := mockApi(ctx, dir)
	defer mock.Close()
	u, _ := url.Parse(mock.URL)

	s = time.Now().UTC()
	hr, _ := getter.GetUrl(u)
	e = time.Now().UTC()
	dur = e.Sub(s)

	slog.Warn("api call duration", slog.Float64("seconds", dur.Seconds()))

	slog.Warn("SeededApiCallOnly",
		slog.Int("N", N),
		slog.String("u", u.String()))

	slog.Warn("end")

	_, bytes := convert.Stringify(hr)
	response := resp.New()
	convert.Unmarshal(bytes, response)

	counts := response.Metadata["counters"].(map[string]interface{})
	all := counts["totals"].(map[string]interface{})

	count := int(all["count"].(float64))
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", all)
	}

	// -- call other api urls and check response
	list := []string{"?archived=true", "?archived=true&team=foo", "?team=foo"}
	for _, l := range list {
		s = time.Now().UTC()
		call := u.String() + l
		ur, _ := url.Parse(call)
		hr, _ = getter.GetUrl(ur)

		e = time.Now().UTC()
		dur = e.Sub(s)
		slog.Warn("api call duration", slog.Float64("seconds", dur.Seconds()), slog.String("url", ur.String()))
		if hr.StatusCode != http.StatusOK {
			t.Errorf("api call failed")
		}
	}

}

func mockApi(ctx context.Context, dir string) *httptest.Server {
	dbF := filepath.Join(dir, "github_standards.db")
	mux := testhelpers.Mux()
	funcs := github_standards.Handlers(ctx, mux, dbF)
	list := funcs["list"]

	return testhelpers.MockServer(list, "warn")
}

func testSeed(dbName string, N int) *seeder.Seed {
	lines := []string{}
	owner := fake.String(12)
	for x := 0; x < N; x++ {
		id := N + x
		g := ghs.Fake(&id, &owner)
		if x == 0 {
			lines = append(lines, g.CSVHead())
		}
		lines = append(lines, g.ToCSV())
	}

	return &seeder.Seed{
		Label:  "test",
		Table:  "github_standards",
		DB:     dbName,
		Schema: "../../../datastore/github_standards/github_standards.sql",
		Source: []string{},
		Dummy:  lines,
	}
}
