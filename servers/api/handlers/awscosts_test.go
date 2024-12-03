package handlers_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/strutils"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
	"github.com/ministryofjustice/opg-reports/servers/inputs"
)

// TestApiHandlersAwsCostsDetailedHandler
func TestApiHandlersAwsCostsDetailedHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsSumFullDetailsResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
		expected map[string]string = map[string]string{}
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](100)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit
		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.RequiredGroupedDateRangeUnitInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	response, err = handlers.ApiAwsCostsSumFullDetailsHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsCostsSumFullDetailsOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsSumFullDetailsOperationID, response.Body.Operation)
	}

	// work out expected sum values for each month
	for _, c := range costs {
		var ym = dateutils.Reformat(c.Date, dateformats.YM)
		var key = fmt.Sprintf("%s.%s.%s.%s.%s.%s", c.Service, c.Region, c.AwsAccount.Number, c.AwsAccount.Environment, c.Unit.Name, ym)
		if _, ok := expected[key]; !ok {
			expected[key] = "0.0"
		}
		if c.Service != "Tax" {
			expected[key] = strutils.Adds(expected[key], c.Cost)
		}
	}

	// now compare and make sure the actual matches the expected value
	for _, res := range response.Body.Result {
		var ym = res.Date
		var key = fmt.Sprintf("%s.%s.%s.%s.%s.%s", res.Service, res.Region, res.AwsAccountNumber, res.AwsAccountEnvironment, res.UnitName, ym)
		var expect = strutils.FloatF(expected[key])
		var actual = strutils.FloatF(res.Cost)
		if expect != actual {
			t.Errorf("costs for [%s] did not match - expected [%s] actual [%v]", key, expect, actual)
		}
	}

}

// TestApiHandlersAwsCostsSumPerUnitEnvHandler
func TestApiHandlersAwsCostsSumPerUnitEnvHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsSumPerUnitEnvResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
		expected map[string]string = map[string]string{}
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](100)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit
		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.RequiredGroupedDateRangeInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	response, err = handlers.ApiAwsCostsSumPerUnitEnvHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsCostsSumPerUnitEnvOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsSumPerUnitEnvOperationID, response.Body.Operation)
	}

	// work out expected sum values for each month
	for _, c := range costs {
		var ym = dateutils.Reformat(c.Date, dateformats.YM)
		var key = fmt.Sprintf("%s.%s.%s", ym, c.Unit.Name, c.AwsAccount.Environment)
		if _, ok := expected[key]; !ok {
			expected[key] = "0.0"
		}
		if c.Service != "Tax" {
			expected[key] = strutils.Adds(expected[key], c.Cost)
		}
	}

	// now compare and make sure the actual matches the expected value
	for _, res := range response.Body.Result {
		var ym = res.Date
		var key = fmt.Sprintf("%s.%s.%s", ym, res.UnitName, res.AwsAccountEnvironment)
		var expect = strutils.FloatF(expected[key])
		var actual = strutils.FloatF(res.Cost)
		if expect != actual {
			t.Errorf("costs for [%s] did not match - expected [%s] actual [%v]", key, expect, actual)
		}
	}

}

// TestApiHandlersAwsCostsSumPerUnitHandler
func TestApiHandlersAwsCostsSumPerUnitHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsSumPerUnitResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
		expected map[string]string = map[string]string{}
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](100)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit
		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.RequiredGroupedDateRangeInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	response, err = handlers.ApiAwsCostsSumPerUnitHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsCostsSumPerUnitOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsSumPerUnitOperationID, response.Body.Operation)
	}

	// work out expected sum values for each month
	for _, c := range costs {
		var ym = dateutils.Reformat(c.Date, dateformats.YM)
		var key = fmt.Sprintf("%s.%s", ym, c.Unit.Name)
		if _, ok := expected[key]; !ok {
			expected[key] = "0.0"
		}
		if c.Service != "Tax" {
			expected[key] = strutils.Adds(expected[key], c.Cost)
		}
	}

	// now compare and make sure the actual matches the expected value
	for _, res := range response.Body.Result {
		var ym = res.Date
		var key = fmt.Sprintf("%s.%s", ym, res.UnitName)
		var expect = strutils.FloatF(expected[key])
		var actual = strutils.FloatF(res.Cost)
		if expect != actual {
			t.Errorf("costs for [%s] did not match - expected [%s] actual [%v]", key, expect, actual)
		}
	}

}

// TestApiHandlersAwsCostsSumHandler
func TestApiHandlersAwsCostsSumHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsSumResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
		expected map[string]string = map[string]string{}
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](100)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit
		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.RequiredGroupedDateRangeUnitInput{
		Version:   "v1",
		Interval:  "month",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	in.Resolve(nil)
	response, err = handlers.ApiAwsCostsSumHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsCostsSumOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsSumOperationID, response.Body.Operation)
	}

	// work out expected sum values for each month
	for _, c := range costs {
		var ym = dateutils.Reformat(c.Date, dateformats.YM)
		if _, ok := expected[ym]; !ok {
			expected[ym] = "0.0"
		}
		if c.Service != "Tax" {
			expected[ym] = strutils.Adds(expected[ym], c.Cost)
		}
	}
	// now compare and make sure the actual matches the expected value
	for _, res := range response.Body.Result {
		var ym = res.Date
		var expect = strutils.FloatF(expected[ym])
		var actual = strutils.FloatF(res.Cost)
		if expect != actual {
			t.Errorf("costs for [%s] did not match - expected [%s] actual [%v]", ym, expect, actual)
		}
	}

}

// TestApiHandlersAwsCostsListHandler inserts dummy data into the db, calls the api handler for
// lists and checks the count of results returned is the same as number inserted
func TestApiHandlersAwsCostsListHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsListResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](50)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit

		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.DateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	response, err = handlers.ApiAwsCostsListHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check data returned
	// check the response info
	if handlers.AwsCostsListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsListOperationID, response.Body.Operation)
	}

	if len(costs) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected at least [%d] actual [%v]", len(costs), len(response.Body.Result))
	}
}

func TestApiHandlersAwsCostsTaxesHandler(t *testing.T) {
	var (
		err     error
		adaptor dbs.Adaptor
		// response *handlers.AwsCostsTaxesResponse
		dir string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](50)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit

		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.RequiredGroupedDateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
		Interval:  "month",
	}
	in.Resolve(nil)

	_, err = handlers.ApiAwsCostsTaxesHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

}

// TestApiHandlersAwsCostsTotalHandler
func TestApiHandlersAwsCostsTotalHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.AwsCostsTotalResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		costs    []*models.AwsCost
		units    []*models.Unit
		accounts []*models.AwsAccount
	)
	fakerextras.AddProviders()
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	// boot strap the db
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	costs = fakermany.Fake[*models.AwsCost](5)

	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
	// join costs to accounts & unit
	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var unit = acc.Unit

		cost.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		cost.Unit = unit
	}

	// seed the cost data
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// query the api handler
	in := &inputs.DateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	response, err = handlers.ApiAwsCostsTotalHandler(ctx, in)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsCostsTotalOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsCostsTotalOperationID, response.Body.Operation)
	}
	actual := "0.0"
	expected := "0.0"
	for _, cost := range costs {
		if cost.Service != "Tax" {
			expected = strutils.Adds(expected, cost.Cost)
		}
	}

	for _, row := range response.Body.Result {
		actual = strutils.Adds(actual, row.Cost)
	}

	actual = strutils.FloatF(actual)
	expected = strutils.FloatF(expected)
	if expected != actual {
		t.Errorf("totals did not match up - expected [%s] actual [%s]", expected, actual)
	}
}
