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
	"github.com/ministryofjustice/opg-reports/servers/api/inputs"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
)

func TestApiHandlersGitHubRepositoryStandardsListHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.GitHubRepositoryStandardsListResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile    string          = filepath.Join(dir, "test.db")
		ctxKey    string          = lib.CTX_DB_KEY
		ctx       context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		repos     []*models.GitHubRepository
		standards []*models.GitHubRepositoryStandard
	)
	fakerextras.AddProviders()

	repos = fakermany.Fake[*models.GitHubRepository](5)
	standards = fakermany.Fake[*models.GitHubRepositoryStandard](5)

	for i, st := range standards {
		st.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repos[i])
		st.GitHubRepositoryFullName = st.GitHubRepository.FullName
	}

	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// seed the teams and units
	_, err = seed.GitHubStandards(ctx, adaptor, standards)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// should return everything
	response, err = handlers.ApiGitHubRepositoryStandardsListHandler(ctx, &inputs.VersionUnitInput{
		Version: "v1",
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.GitHubRepositoryStandardsListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.GitHubRepositoryStandardsListOperationID, response.Body.Operation)
	}
	// check the number of results
	if len(standards) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(standards), len(response.Body.Result))
	}

}
