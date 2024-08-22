package github_standards_test

import (
	"context"
	"log"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema string = "../../../datastore/github_standards/github_standards.sql"

func TestServersFrontGithubStandards(t *testing.T) {
	logger.LogSetup()

	//--- spin up an api
	// seed
	ctx := context.TODO()
	N := 10
	dir := t.TempDir()
	dbF := filepath.Join(dir, "ghs.db")
	schemaF := filepath.Join(dir, "ghs.sql")
	dataF := filepath.Join(dir, "dummy.json")
	testhelpers.CopyFile(realSchema, schemaF)
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "github_standards", N)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	// set mock
	github_standards.SetDBPath(dbF)
	github_standards.SetCtx(ctx)
	mock := testhelpers.MockServer(github_standards.ListHandler, "warn")
	defer mock.Close()

	// -- mimic to call local
}
