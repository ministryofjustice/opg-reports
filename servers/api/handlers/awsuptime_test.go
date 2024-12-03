package handlers_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

func TestApiHandlersAwsUptimeList(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.AwsUptimeListResponse
		dir      string = t.TempDir()
		// dir       string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		uptime   []*models.AwsUptime
		accounts []*models.AwsAccount
		units    []*models.Unit
	)

	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	uptime = fakermany.Fake[*models.AwsUptime](20)

	for _, up := range uptime {
		var acc = fakerextras.Choice(accounts)
		var unt = fakerextras.Choice(units)
		up.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		up.Unit = (*models.UnitForeignKey)(unt)
	}

	_, err = seed.AwsUptime(ctx, adaptor, uptime)
	if err != nil {
		t.Fatalf(err.Error())
	}

	in := &inout.DateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}

	// should return everything as we are only using 1 unit
	response, err = handlers.ApiAwsUptimeListHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsUptimeListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsUptimeListOperationID, response.Body.Operation)
	}
	if len(uptime) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected at least [%d] actual [%v]", len(uptime), len(response.Body.Result))
	}
}

func TestApiHandlersAwsUptimeAverages(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.AwsUptimeAveragesResponse
		dir      string = t.TempDir()
		// dir      string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		uptime   []*models.AwsUptime
		accounts []*models.AwsAccount
		units    []*models.Unit
	)

	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	uptime = fakermany.Fake[*models.AwsUptime](200)

	for _, up := range uptime {
		var acc = fakerextras.Choice(accounts)
		var unt = fakerextras.Choice(units)
		up.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		up.Unit = (*models.UnitForeignKey)(unt)
	}

	_, err = seed.AwsUptime(ctx, adaptor, uptime)
	if err != nil {
		t.Fatalf(err.Error())
	}

	in := &inout.RequiredGroupedDateRangeUnitInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	// should return everything as we are only using 1 unit
	response, err = handlers.ApiAwsUptimeAveragesHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsUptimeAveragesOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsUptimeAveragesOperationID, response.Body.Operation)
	}

}

func TestApiHandlersAwsUptimeAveragesPerUnit(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.AwsUptimeAveragesPerUnitResponse
		dir      string = t.TempDir()
		// dir      string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		uptime   []*models.AwsUptime
		accounts []*models.AwsAccount
		units    []*models.Unit
	)

	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	uptime = fakermany.Fake[*models.AwsUptime](200)

	for _, up := range uptime {
		var acc = fakerextras.Choice(accounts)
		var unt = fakerextras.Choice(units)
		up.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		up.Unit = (*models.UnitForeignKey)(unt)
	}

	_, err = seed.AwsUptime(ctx, adaptor, uptime)
	if err != nil {
		t.Fatalf(err.Error())
	}

	in := &inout.RequiredGroupedDateRangeInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	// should return everything as we are only using 1 unit
	response, err = handlers.ApiAwsUptimeAveragesPerUnitHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	total := 0
	for _, r := range response.Body.Result {
		total += r.Count
	}

	// check the response info
	if handlers.AwsUptimeAveragesPerUnitOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsUptimeAveragesPerUnitOperationID, response.Body.Operation)
	}

	// the total sum of the counts should match the number creates
	if len(uptime) != total {
		t.Errorf("error with number of results - expected at least [%d] actual [%v]", len(uptime), total)
	}
}
