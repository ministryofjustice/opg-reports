package accountimport

import (
	"context"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/files"
	"opg-reports/report/package/logger"
	"path/filepath"
	"testing"
)

func TestAccountImportWithoutMock(t *testing.T) {
	var (
		err     error
		ctx     context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir     string          = t.TempDir()
		srcfile string          = filepath.Join(dir, "aws.accounts.json")
		dbpath  string          = filepath.Join(dir, "test-accounts-import.db")
	)
	// generate some teams
	data := []*Model{
		{
			ID:          "A001",
			Name:        "Test 01",
			Label:       "test",
			Environment: "development",
			TeamName:    "team-dev",
		},
		{
			ID:          "A002",
			Name:        "Test 02",
			Label:       "test",
			Environment: "production",
			TeamName:    "team-production",
		},
	}
	err = files.WriteAsJSON(ctx, srcfile, data)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	// db migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     dbpath,
		Driver: "sqlite3",
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
