package aws_costs_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema string = "../../../datastore/aws_costs/aws_costs.sql"

func TestServersApiAwsCostsApiYtd(t *testing.T) {

	logger.LogSetup()
	ctx := context.TODO()
	N := 5000
	dir := t.TempDir()

	dbF := filepath.Join(dir, "awsc.db")
	schemaF := filepath.Join(dir, "awsc.sql")
	dataF := filepath.Join(dir, "dummy.json")

	testhelpers.CopyFile(realSchema, schemaF)
	tick := testhelpers.T()
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "aws_costs", N)
	tick.Stop()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	slog.Warn("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	aws_costs.SetDBPath(dbF)
	aws_costs.SetCtx(ctx)
	// -- setup mock api
	mock := testhelpers.MockServer(aws_costs.YtdHandler, "warn")
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

	slog.Warn("api call duration", slog.String("seconds", tick.Seconds()), slog.String("u", u.String()))

	// -- check values of the response
	_, bytes := convert.Stringify(hr)
	response, _ := convert.Unmarshal[*resp.Response](bytes)

	counts := response.Metadata["counters"].(map[string]interface{})
	all := counts["totals"].(map[string]interface{})
	count := int(all["count"].(float64))
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", all)
	}

}

func TestServersApiAwsCostsApiMonthlyTax(t *testing.T) {

	logger.LogSetup()
	ctx := context.TODO()
	N := 5000
	dir := t.TempDir()

	dbF := filepath.Join(dir, "awsc.db")
	schemaF := filepath.Join(dir, "awsc.sql")
	dataF := filepath.Join(dir, "dummy.json")

	testhelpers.CopyFile(realSchema, schemaF)
	tick := testhelpers.T()
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "aws_costs", N)
	tick.Stop()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	slog.Warn("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	aws_costs.SetDBPath(dbF)
	aws_costs.SetCtx(ctx)
	// -- setup mock api
	mock := testhelpers.MockServer(aws_costs.MonthlyTaxHandler, "warn")
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

	slog.Warn("api call duration", slog.String("seconds", tick.Seconds()), slog.String("u", u.String()))

	// -- check values of the response
	_, bytes := convert.Stringify(hr)
	response, _ := convert.Unmarshal[*resp.Response](bytes)

	counts := response.Metadata["counters"].(map[string]interface{})
	all := counts["totals"].(map[string]interface{})
	count := int(all["count"].(float64))
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", all)
	}

}

func TestServersApiAwsCostsApiStandard(t *testing.T) {

	logger.LogSetup()
	ctx := context.TODO()
	N := 10000
	dir := t.TempDir()

	dbF := filepath.Join(dir, "awsc.db")
	schemaF := filepath.Join(dir, "awsc.sql")
	dataF := filepath.Join(dir, "dummy.json")

	testhelpers.CopyFile(realSchema, schemaF)
	tick := testhelpers.T()
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "aws_costs", N)
	tick.Stop()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	slog.Warn("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	aws_costs.SetDBPath(dbF)
	aws_costs.SetCtx(ctx)
	// -- setup mock api
	mock := testhelpers.MockServer(aws_costs.StandardHandler, "warn")
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

	slog.Warn("api call duration", slog.String("seconds", tick.Seconds()), slog.String("u", u.String()))

	// -- check values of the response
	_, bytes := convert.Stringify(hr)
	response, _ := convert.Unmarshal[*resp.Response](bytes)

	counts := response.Metadata["counters"].(map[string]interface{})
	all := counts["totals"].(map[string]interface{})
	count := int(all["count"].(float64))
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", all)
	}

	// -- call with options
	list := []string{"?group=unit-env", "?group=detailed", "?group=unit"}
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

		slog.Warn("api call duration", slog.String("seconds", tick.Seconds()), slog.String("url", ur.String()))
		if hr.StatusCode != http.StatusOK {
			t.Errorf("api call failed")
		}
	}

}
