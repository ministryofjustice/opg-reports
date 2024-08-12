package github_standards_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/seeder/github_standards_seed"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestServersApiGithubStandardsArchivedPerf(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	N := 500000

	// dir := testhelpers.Dir()
	// slog.Info(dir)
	// defer os.RemoveAll(dir)

	db, err := newDb(ctx, dir)
	defer db.Close()
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	// time seeding
	s := time.Now().UTC()
	q := seedDB(ctx, db, N)

	e := time.Now().UTC()
	dur := e.Sub(s)
	slog.Warn("seed duration", slog.Float64("seconds", dur.Seconds()))

	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	s = time.Now().UTC()
	q.ArchivedFilter(ctx, 1)
	e = time.Now().UTC()
	dur = e.Sub(s)
	slog.Warn("archived filter duration", slog.Float64("seconds", dur.Seconds()))

	// // -- generate a dummy database using seeder
	// // 	copy over schema from normal location to tmp dir
	// // 	WARNING - presumes running from root dir / make file
	// // 		for the path
	// path := "../../../datastore/github_standards"
	// scF := dir + "ghs.sql"
	// dbF := dir + "ghs.db"
	// defer os.Remove(scF)
	// defer os.Remove(dbF)
	// testhelpers.CopyFile(path+"/schema.sql", scF)

	// db, _ := github_standards_seed.NewDb(ctx, dbF, scF)
	// q := github_standards_seed.Seed(ctx, db, N)

	// l, _ := q.Count(ctx)
	// if l != int64(create) {
	// 	t.Errorf("records did not create properly")
	// }

	// // -- now test calling / binding to the endpoint
	// mux := testhelpers.Mux()
	// funcs := github_standards.Handlers(ctx, mux, dbF)
	// list := funcs["list"]

	// mock := testhelpers.MockServer(list, "warn")
	// defer mock.Close()
	// u, _ := url.Parse(mock.URL)

	// hr, _ := getter.GetUrl(u)

	// s, bytes := convert.Stringify(hr)
	// response := resp.New()
	// convert.Unmarshal(bytes, response)

	// if len(response.Result) != create {
	// 	t.Errorf("incorrect amount returned")
	// 	fmt.Println(s)
	// }
}

// Benchmark to see how long it takes for the mock api to
// return the data we want
func BenchmarkServerApiGithubStandardsSeededApiCallOnly(b *testing.B) {
	logger.LogSetup()
	num := b.N
	slog.Warn("SeededApiCallOnly", slog.Int("N", num))

	ctx := context.Background()
	dir := b.TempDir()
	b.StopTimer()
	b.ResetTimer()

	db, err := newDb(ctx, dir)
	if err != nil {
		b.Errorf("error with db:" + err.Error())
	}
	defer os.RemoveAll(dir)
	defer db.Close()
	seedDB(ctx, db, num)
	mock := mockApi(ctx, dir)
	defer mock.Close()
	u, _ := url.Parse(mock.URL)

	b.StopTimer()
	b.StartTimer()
	getter.GetUrl(u)

	s := b.Elapsed().Seconds()
	b.StopTimer()

	slog.Warn("SeededApiCallOnly",
		slog.Int("N", num),
		slog.String("u", u.String()),
		slog.String("dir", dir),
		slog.Float64("seconds", s),
	)

}

// func BenchmarkServerApiGithubStandardsAll(b *testing.B) {
// 	logger.LogSetup()
// 	ctx := context.Background()
// 	b.StopTimer()
// 	b.StartTimer()

// 	create := 500000

// 	// -- time hwo long the db creation takes
// 	b.ResetTimer()
// 	b.StopTimer()
// 	b.StartTimer()

// 	slog.Info("creating db: " + dbF)

// 	db.Ping()
// 	slog.Info("seeding db", slog.Int("count", create))
// 	github_standards_seed.Seed(ctx, db, create)
// 	defer db.Close()

// 	b.StopTimer()

// 	mux := testhelpers.Mux()
// 	funcs := github_standards.Handlers(ctx, mux, dbF)
// 	list := funcs["list"]

// 	mock := testhelpers.MockServer(list)
// 	defer mock.Close()
// 	uri1, _ := url.Parse(mock.URL)
// 	uri2, _ := url.Parse(mock.URL + "?archived=false")

// 	set := []*url.URL{uri1, uri2}
// 	for _, u := range set {

// 		b.StartTimer()
// 		hr, err := getter.GetUrl(u)
// 		b.StopTimer()

// 		b.StartTimer()
// 		if err != nil || hr.StatusCode != http.StatusOK {
// 			b.Fatal("request failed")
// 			fmt.Println(err)
// 			fmt.Println(hr.StatusCode)
// 		} else {
// 			_, bytes := convert.Stringify(hr)
// 			response := resp.New()
// 			// m := response.Metadata
// 			convert.Unmarshal(bytes, response)
// 			// slog.Info("duration", slog.Float64("seconds", response.Timer.Duration.Seconds()))
// 			// slog.Info("metadata", slog.String("total", fmt.Sprintf("%+v\n", m)))
// 		}
// 		b.StopTimer()

// 	}

// }

func newDb(ctx context.Context, dir string) (*sql.DB, error) {
	path := "../../../datastore/github_standards"
	scF := filepath.Join(dir, "ghs.sql")
	dbF := filepath.Join(dir, "ghs.db")
	s := filepath.Join(path, "schema.sql")
	testhelpers.CopyFile(s, scF)

	return github_standards_seed.NewDb(ctx, dbF, scF)
}

func seedDB(ctx context.Context, db *sql.DB, num int) *ghs.Queries {
	return github_standards_seed.Seed(ctx, db, num)
}

func mockApi(ctx context.Context, dir string) *httptest.Server {
	dbF := filepath.Join(dir, "ghs.db")
	mux := testhelpers.Mux()
	funcs := github_standards.Handlers(ctx, mux, dbF)
	list := funcs["list"]

	return testhelpers.MockServer(list, "warn")
}
