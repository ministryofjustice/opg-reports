package seeder

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema = "../../../datastore/github_standards/github_standards.sql"

// Seeder to make a new db using a real schema from existing (generated)
// data files with dummy data
func TestSeedSeederNewDbWithSchemaFromDataFiles(t *testing.T) {
	logger.LogSetup()
	ctx := context.TODO()
	n := 100
	dir := t.TempDir()

	dbF := filepath.Join(dir, "ghs.db")
	schemaF := filepath.Join(dir, "ghs.sql")
	dataF := filepath.Join(dir, "dummy.json")
	// copy over the real schema
	testhelpers.CopyFile(realSchema, schemaF)
	// -- create the dummy files
	owner := fake.String(12)
	list := []*ghs.GithubStandard{}
	for x := 0; x < n; x++ {
		list = append(list, ghs.Fake(nil, &owner))
	}

	// write to dummy file
	content, err := convert.Marshals(list)
	if err != nil {
		t.Errorf("error with marshaling:" + err.Error())
	}
	os.WriteFile(dataF, content, os.ModePerm)

	// -- now seed the database
	tick := testhelpers.T()

	db, err := Seed(ctx, dbF, schemaF, dataF, "github_standards", n)
	defer db.Close()
	tick.Stop()

	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	// -- check counts
	slog.Warn("seed duration", slog.String("seconds", tick.Seconds()))
	q := ghs.New(db)
	l, err := q.Count(ctx)
	slog.Warn("count check", slog.Int64("found", l), slog.Int("n", n))
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	if l != int64(n) {
		t.Errorf("records did not create properly: [%d] [%d]", n, l)
	}

}

// Seeder to make a new database from real schema into a temp dir
// with directly generated dummy data
func TestSeedSeederNewDbWithSchemaNoDataFiles(t *testing.T) {
	logger.LogSetup()
	ctx := context.TODO()
	n := 100
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
	slog.Info("ending seed")
	tick.Stop()
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	slog.Warn("seed duration", slog.String("seconds", tick.Seconds()))
	q := ghs.New(db)
	l, err := q.Count(ctx)
	slog.Warn("count check", slog.Int64("found", l), slog.Int("n", n))
	if err != nil {
		t.Errorf("error with db:" + err.Error())
	}

	if l != int64(n) {
		t.Errorf("records did not create properly: [%d] [%d]", n, l)
	}

}

// Benchmark the perofmrance of injecting 1000 recrods in N times
// to track duration
func BenchmarkSeedSeederNewDbWithSchemaNoDataFiles1k(b *testing.B) {
	logger.LogSetup()
	ctx := context.TODO()
	n := 1000
	tick := testhelpers.T()
	for i := 0; i < b.N; i++ {
		dir := testhelpers.Dir()
		defer os.RemoveAll(dir)

		dbF := filepath.Join(dir, "ghs.db")
		schemaF := filepath.Join(dir, "ghs.sql")
		dataF := filepath.Join(dir, "*.json")
		// copy over the real schema
		testhelpers.CopyFile(realSchema, schemaF)
		Seed(ctx, dbF, schemaF, dataF, "github_standards", n)

	}
	slog.Warn("bench", slog.Int("N", b.N*n), slog.Int("b.N", b.N), slog.String("seconds", tick.Stop().Seconds()))

}
