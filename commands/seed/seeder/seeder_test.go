package seeder

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema = "../../../datastore/github_standards/github_standards.sql"

func TestSeedSeederNoDbSchemaExistsNoDataFile(t *testing.T) {
	logger.LogSetup()
	//
	ctx := context.TODO()
	n := 500000
	dir := t.TempDir()

	dbF := filepath.Join(dir, "ghs.db")
	schemaF := filepath.Join(dir, "ghs.sql")
	dataF := filepath.Join(dir, "*.json")
	// copy over the real schema
	testhelpers.CopyFile(realSchema, schemaF)
	// -- seed and check the time
	tick := testhelpers.T()
	slog.Info("starting seed")
	db, err := Seed(ctx, dbF, schemaF, dataF, "github_standards", n)
	defer db.Close()
	slog.Warn("ending seed")
	tick.Stop()
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	slog.Info("seed duration", slog.String("seconds", tick.Seconds()))
	q := ghs.New(db)
	l, err := q.Count(ctx)
	slog.Info("count check", slog.Int64("found", l), slog.Int("n", n))
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	if l != int64(n) {
		t.Errorf("records did not create properly: [%d] [%d]", n, l)
	}

}
