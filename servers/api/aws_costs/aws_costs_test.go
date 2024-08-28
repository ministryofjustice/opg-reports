package aws_costs_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
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
	slog.Debug("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	server := api.New(ctx, dbF)
	handler := api.Wrap(server, aws_costs.YtdHandler)
	// -- setup mock api
	mock := testhelpers.MockServer(handler, "warn")
	defer mock.Close()

	// -- call the api - time its duration
	hr, err := httphandler.Get("", "", mock.URL)

	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Debug("api call duration", slog.Float64("seconds", hr.Duration), slog.String("url", mock.URL))

	// -- check values of the response
	_, bytes := convert.Stringify(hr.Response)
	response, _ := convert.Unmarshal[*aws_costs.ApiResponse](bytes)

	counts := response.Counters
	count := counts.Totals.Count
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", counts)
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
	slog.Debug("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	server := api.New(ctx, dbF)
	handler := api.Wrap(server, aws_costs.MonthlyTaxHandler)
	// -- setup mock api
	mock := testhelpers.MockServer(handler, "warn")
	defer mock.Close()
	// -- call the api - time its duration

	hr, err := httphandler.Get("", "", mock.URL)

	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Debug("api call duration", slog.Float64("seconds", hr.Duration), slog.String("url", mock.URL))

	// -- check values of the response
	_, bytes := convert.Stringify(hr.Response)
	response, _ := convert.Unmarshal[*aws_costs.ApiResponse](bytes)

	counts := response.Counters
	count := counts.Totals.Count
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", counts)
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
	slog.Debug("seed duration", slog.String("seconds", tick.Seconds()))

	// check the count of records
	q := awsc.New(db)
	defer q.Close()
	l, _ := q.Count(ctx)
	if l != int64(N) {
		t.Errorf("records did not create properly: [%d] [%d]", N, l)
	}

	// -- set db and context
	server := api.New(ctx, dbF)
	handler := api.Wrap(server, aws_costs.StandardHandler)
	// -- setup mock api
	mock := testhelpers.MockServer(handler, "warn")
	defer mock.Close()

	// -- call the api - time its duration
	hr, err := httphandler.Get("", "", mock.URL)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}

	slog.Debug("api call duration", slog.Float64("seconds", hr.Duration), slog.String("url", mock.URL))

	// -- check values of the response
	_, bytes := convert.Stringify(hr.Response)
	response, _ := convert.Unmarshal[*aws_costs.ApiResponse](bytes)

	counts := response.Counters
	count := counts.Totals.Count
	if count != N {
		t.Errorf("total number of rows dont match")
		fmt.Printf("%+v\n", counts)
	}

	// -- call with options
	list := []string{"?group=unit-env", "?group=detailed", "?group=unit"}
	for _, l := range list {

		hr, err = httphandler.Get("", "", mock.URL+l)
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
