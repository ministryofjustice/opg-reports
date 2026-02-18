package teamimport

import (
	"context"
	"opg-reports/report/internal/global"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/files"
	"opg-reports/report/package/logger"
	"path/filepath"
	"testing"
)

func TestTeamImportWithoutMock(t *testing.T) {

	var (
		err     error
		ctx     context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir     string          = t.TempDir()
		mfile   string          = filepath.Join(dir, "migrate.json")
		srcfile string          = filepath.Join(dir, "teams.json")
		dbpath  string          = filepath.Join(dir, "test-costs-import.db")
	)
	// generate some teams
	teams := []*TeamModel{
		{Name: "team-a"}, {Name: "team-b"}, {Name: "team-c"}, {Name: "team-d"}, {Name: "team-e"},
	}
	err = files.WriteAsJSON(ctx, srcfile, teams)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// db migrations
	err = global.MigrateAll(ctx, &migrations.Args{
		DB:            dbpath,
		Driver:        "sqlite3",
		MigrationFile: mfile,
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	err = Import(ctx, &Args{
		DB:      dbpath,
		Driver:  "sqlite3",
		SrcFile: srcfile,
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

}
