package github_standards_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/seeder/github_standards_seed"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestServersApiGithubStandards(t *testing.T) {
	ctx := context.Background()
	dir := testhelpers.Dir()
	create := 5
	defer os.RemoveAll(dir)
	// -- generate a dummy database using seeder
	// 	copy over schema from normal location to tmp dir
	// 	WARNING - presumes running from root dir / make file
	// 		for the path
	path := "../../../datastore/github_standards"
	scF := dir + "ghs.sql"
	dbF := dir + "ghs.db"
	defer os.Remove(scF)
	defer os.Remove(dbF)
	testhelpers.CopyFile(path+"/schema.sql", scF)

	db := github_standards_seed.NewDb(ctx, dbF, scF)
	q := github_standards_seed.Seed(ctx, db, create)

	l, _ := q.Count(ctx)
	if l != int64(create) {
		t.Errorf("records did not create properly")
	}

	// -- now test calling / binding to the endpoint
	mux := testhelpers.Mux()
	funcs := github_standards.Handlers(ctx, mux, dbF)
	list := funcs["list"]

	mock := testhelpers.MockServer(list)
	defer mock.Close()
	u, _ := url.Parse(mock.URL)

	hr, _ := getter.GetUrl(u)

	s, bytes := convert.Stringify(hr)
	response := resp.New()
	convert.Unmarshal(bytes, response)

	if len(response.Result) != create {
		t.Errorf("incorrect amount returned")
		fmt.Println(s)
	}
}

func BenchmarkServerApiGithubStandardsAll(b *testing.B) {
	logger.LogSetup()
	ctx := context.Background()
	dir := testhelpers.Dir()
	create := 500000
	defer os.RemoveAll(dir)
	// -- make db
	path := "../../../datastore/github_standards"
	scF := dir + "ghs.sql"
	dbF := dir + "ghs.db"
	defer os.Remove(scF)
	defer os.Remove(dbF)
	testhelpers.CopyFile(path+"/schema.sql", scF)

	slog.Info("creating db: " + dbF)
	db := github_standards_seed.NewDb(ctx, dbF, scF)
	db.Ping()

	slog.Info("seeding db", slog.Int("count", create))
	github_standards_seed.Seed(ctx, db, create)
	defer db.Close()

	mux := testhelpers.Mux()
	funcs := github_standards.Handlers(ctx, mux, dbF)
	list := funcs["list"]

	mock := testhelpers.MockServer(list)
	defer mock.Close()
	uri1, _ := url.Parse(mock.URL)
	uri2, _ := url.Parse(mock.URL + "?archive=true")

	set := []*url.URL{uri1, uri2}
	for _, u := range set {
		b.ResetTimer()
		slog.Info("calling handler", slog.String("url", mock.URL))
		hr, err := getter.GetUrl(u)
		if err != nil || hr.StatusCode != http.StatusOK {
			b.Fatal("request failed")
			fmt.Println(err)
			fmt.Println(hr.StatusCode)
		} else {
			_, bytes := convert.Stringify(hr)
			response := resp.New()
			convert.Unmarshal(bytes, response)
			slog.Info("duration", slog.Float64("seconds", response.Timer.Duration.Seconds()))
		}
	}

}
