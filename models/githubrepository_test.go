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
	_ dbs.Table           = &models.GitHubRepository{}
	_ dbs.CreateableTable = &models.GitHubRepository{}
	_ dbs.Insertable      = &models.GitHubRepository{}
	_ dbs.Row             = &models.GitHubRepository{}
	_ dbs.InsertableRow   = &models.GitHubRepository{}
	_ dbs.Record          = &models.GitHubRepository{}
)

// TestModelsGitHubRepositoryCRUD checks the github team table creation
// and inserting series of fake records works as expected
func TestModelsGitRepositoryCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int                        = 100
		ctx     context.Context            = context.Background()
		dir     string                     = t.TempDir()
		path    string                     = filepath.Join(dir, "test.db")
		units   []*models.GitHubRepository = fakermany.Fake[*models.GitHubRepository](n)
		results []*models.GitHubRepository
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.GitHubRepository{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.GitHubRepository{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.GitHubRepository{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}
