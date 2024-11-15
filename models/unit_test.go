package models_test

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
)

// Interface checks
var (
	_ dbs.Table           = &models.Unit{}
	_ dbs.CreateableTable = &models.Unit{}
	_ dbs.Insertable      = &models.Unit{}
	_ dbs.Row             = &models.Unit{}
	_ dbs.InsertableRow   = &models.Unit{}
	_ dbs.Record          = &models.Unit{}
)

// TestModelsUnitCRUD checks the unit table and inserting series of fake
// records works as expected
func TestModelsUnitCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int             = 4
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		units   []*models.Unit  = fakermany.Fake[*models.Unit](n)
		results []*models.Unit
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.Unit{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.Unit{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}
