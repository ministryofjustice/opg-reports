package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCMDSeed(t *testing.T) {

	var (
		err    error
		dir    string = "./"
		dbPath string = filepath.Join(dir, "test-cmd-seed.db")
	)
	// change location of the db
	os.Setenv("DB_PATH", dbPath)
	setup()

	err = seedFunc(nil, []string{})
	if err != nil {
		t.Errorf("unexpected seeding error: %s", err.Error())
	}

}
