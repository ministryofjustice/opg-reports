package handlers_test

import (
	"context"
	"path/filepath"
	"testing"

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

func TestApiHandlerAwsAccountsList(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.AwsAccountsListResponse
		dir      string = t.TempDir()
		// dir       string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		units    []*models.Unit
		accounts []*models.AwsAccount
	)
	fakerextras.AddProviders()

	units = fakermany.Fake[*models.Unit](5)
	accounts = fakermany.Fake[*models.AwsAccount](5)
	for _, acc := range accounts {
		var u = fakerextras.Choice(units)
		acc.Unit = (*models.UnitForeignKey)(u)
	}
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
	// seed the accounts
	_, err = seed.AwsAccounts(ctx, adaptor, accounts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// should return everything
	in := &inout.VersionUnitInput{
		Version: "v1",
	}
	response, err = handlers.ApiAwsAccountsListHandler(ctx, in)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.AwsAccountsListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.AwsAccountsListOperationID, response.Body.Operation)
	}
	if len(accounts) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(accounts), len(response.Body.Result))
	}
}
