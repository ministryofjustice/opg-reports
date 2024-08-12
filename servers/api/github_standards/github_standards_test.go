package github_standards_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/seeder/github_standards_seed"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestServersApiGithubStandardsArchivedPerfDBOnly(t *testing.T) {
	ctx := context.Background()
	N := 500000

	dir := t.TempDir()
	// dir := testhelpers.Dir()
	slog.Warn("dir:" + dir)

	s := time.Now().UTC()
	db, err := seedDb(dir, N)
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
	res, _ := q.ArchivedFilter(ctx, 1)
	e = time.Now().UTC()
	dur = e.Sub(s)
	slog.Warn("archived filter duration",
		slog.Float64("seconds", dur.Seconds()),
		slog.Int("records", len(res)))

}

func TestServersApiGithubStandardsArchivedPerfApiCallOnly(t *testing.T) {
	ctx := context.Background()
	N := 500000
	// dir := t.TempDir()
	dir := testhelpers.Dir()
	slog.Warn("dir:" + dir)
	// defer os.RemoveAll(dir)

	s := time.Now().UTC()
	db, err := seedDb(dir, N)
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

	mock := mockApi(ctx, dir)
	defer mock.Close()
	u, _ := url.Parse(mock.URL)

	s = time.Now().UTC()

	getter.GetUrl(u)

	e = time.Now().UTC()
	dur = e.Sub(s)

	slog.Warn("api call duration", slog.Float64("seconds", dur.Seconds()))

	slog.Warn("SeededApiCallOnly",
		slog.Int("N", N),
		slog.String("u", u.String()),
		slog.String("dir", dir))
	// hr, _ := getter.GetUrl(u)

	// s, bytes := convert.Stringify(hr)
	// response := resp.New()
	// convert.Unmarshal(bytes, response)

	//	if len(response.Result) != create {
	//		t.Errorf("incorrect amount returned")
	//		fmt.Println(s)
	//	}
}

func seedDb(dir string, num int) (*sql.DB, error) {
	path := "../../../datastore/github_standards"
	target := filepath.Join(dir, "github_standards.sql")
	source := filepath.Join(path, "schema.sql")
	testhelpers.CopyFile(source, target)

	return github_standards_seed.NewSeed(dir, num)

}

func mockApi(ctx context.Context, dir string) *httptest.Server {
	dbF := filepath.Join(dir, "github_standards.db")
	mux := testhelpers.Mux()
	funcs := github_standards.Handlers(ctx, mux, dbF)
	list := funcs["list"]

	return testhelpers.MockServer(list, "warn")
}
