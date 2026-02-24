package seeds

import (
	"context"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"path/filepath"
	"testing"
)

func TestSeedsSeedAll(t *testing.T) {
	var (
		err    error
		ctx    context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    string          = t.TempDir()
		mfile  string          = filepath.Join(dir, "migrate.json")
		dbpath string          = filepath.Join(dir, "test-seeds-all.db")
	)

	res, err := SeedAll(ctx, &Args{
		DB:            dbpath,
		Driver:        "sqlite3",
		MigrationFile: mfile,
	})
	if err != nil {
		t.Errorf("unexpcted error : [%s]", err.Error())
	}

	if len(res.Teams) != len(teamList) {
		t.Errorf("not all teams were generated")
	}
	if len(res.Accounts) <= 0 {
		t.Errorf("no accounts generated")
	}
	if len(res.Costs) < 1000 {
		t.Errorf("not enough costs generated")
	}
	if len(res.Uptime) < 100 {
		t.Errorf("not enough uptime records generated")
	}
	if len(res.Codebases) < 10 {
		t.Errorf("not enough codebase records generated")
	}

	// dump.Now(res)

	// t.FailNow()
}
