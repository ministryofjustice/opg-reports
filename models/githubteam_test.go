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
	_ dbs.Table           = &models.GitHubTeam{}
	_ dbs.CreateableTable = &models.GitHubTeam{}
	_ dbs.Insertable      = &models.GitHubTeam{}
	_ dbs.Row             = &models.GitHubTeam{}
	_ dbs.InsertableRow   = &models.GitHubTeam{}
	_ dbs.Record          = &models.GitHubTeam{}
)

// TestModelsGitHubTeamCRUD checks the github team table creation
// and inserting series of fake records works as expected
func TestModelsGitHubTeamCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int                  = 4
		ctx     context.Context      = context.Background()
		dir     string               = t.TempDir()
		path    string               = filepath.Join(dir, "test.db")
		units   []*models.GitHubTeam = fakermany.Fake[*models.GitHubTeam](n)
		results []*models.GitHubTeam
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.GitHubTeam{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.GitHubTeam{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.GitHubTeam{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}
