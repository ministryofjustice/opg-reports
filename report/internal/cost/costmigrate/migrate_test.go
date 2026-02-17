package costmigrate

import (
	"fmt"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dump"
	"opg-reports/report/package/files"
	"opg-reports/report/package/logger"
	"path/filepath"
	"testing"
)

func TestCostMigration(t *testing.T) {
	var (
		err  error
		ctx  = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir  = t.TempDir()
		opts = &Input{
			Driver:        "sqlite3",
			DB:            filepath.Join(dir, "test-costmigration.db"),
			MigrationFile: filepath.Join(dir, "migrations.json"),
		}
	)

	err = Migrate(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	migrated := []string{}
	err = files.ReadJSON(ctx, opts.MigrationFile, &migrated)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	if len(migrated) != len(allMigrations) {
		t.Errorf("not all migrations ran.")
		fmt.Printf("all:\n%v\nmigrated:\n%v\v", dump.Any(allMigrations), dump.Any(migrated))
	}

}
