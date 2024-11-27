package handlers_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/inputs"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
)

func TestApiHandlersUnitsListHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.UnitsListResponse
		dir      string          = t.TempDir()
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		units    []*models.Unit  = []*models.Unit{}
		inserted []*models.Unit  = []*models.Unit{}
	)
	units = fakermany.Fake[*models.Unit](5)
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// insert the dummy units
	inserted, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// check lengths
	if len(units) != len(inserted) {
		t.Errorf("error inserting - expected [%d] actual [%v]", len(units), len(inserted))
	}

	response, err = handlers.ApiUnitsListHandler(ctx, &inputs.VersionInput{
		Version: "v1",
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.UnitListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.UnitListOperationID, response.Body.Operation)
	}
	// check the number of results
	if len(units) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(units), len(response.Body.Result))
	}

}
