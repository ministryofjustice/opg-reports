package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCMDMigrate(t *testing.T) {

	var (
		err    error
		dir    string = t.TempDir()
		dbPath string = filepath.Join(dir, "test-cmd-migrate.db")
	)
	// change location of the db
	os.Setenv("DB_PATH", dbPath)
	setup()

	err = migrateFunc(nil, []string{})
	if err != nil {
		t.Errorf("unexpected migration error: %s", err.Error())
	}

}
