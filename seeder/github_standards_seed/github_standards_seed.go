package github_standards_seed

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/fake"
)

type Seed struct {
	BaseDir    string
	Count      int
	CsvFile    string
	DbFile     string
	SchemaFile string
}

func NewSeed(dir string, count int) (db *sql.DB, err error) {
	seed := &Seed{BaseDir: dir, Count: count}

	// -- file name
	seed.CsvFile = filepath.Join(dir, "github_standards.csv")
	seed.SchemaFile = filepath.Join(dir, "github_standards.sql")
	seed.DbFile = filepath.Join(dir, "github_standards.db")

	os.MkdirAll(seed.BaseDir, os.ModePerm)

	// -- db
	os.WriteFile(seed.DbFile, []byte(""), os.ModePerm)
	// -- connect to db
	conn := consts.SQL_CONNECTION_PARAMS
	db, err = sql.Open("sqlite3", seed.DbFile+conn)
	if err != nil {
		slog.Error("error opening: " + err.Error())
		return
	}
	// -- import schema
	schema, err := os.ReadFile(seed.SchemaFile)
	if err != nil {
		slog.Error("error with schema: " + err.Error())
		return
	} else if _, err = db.Exec(string(schema)); err != nil {
		return
	}

	// -- create dummy data into csv file
	os.WriteFile(seed.CsvFile, []byte(""), os.ModePerm)
	f, err := os.OpenFile(seed.CsvFile, os.O_APPEND|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		slog.Error("error reading csv: "+err.Error(), slog.String("file", seed.CsvFile))
		return
	}
	// -- generate a csv
	owner := fake.String(12)
	for x := 0; x < seed.Count; x++ {
		id := 1000 + x
		g := ghs.Fake(&id, &owner)

		if x == 0 {
			f.WriteString(g.CSVHead())
		}

		line := g.ToCSV()
		f.WriteString(line)
	}
	// -- import csv
	cmd := exec.Command("bash", "-c", "sqlite3", seed.DbFile, "--csv", fmt.Sprintf(".import --skip 1 %s github_standards", seed.CsvFile))
	err = cmd.Run()
	cmd = exec.Command("sqlite3", seed.DbFile, "--csv", fmt.Sprintf(".import --skip 1 %s github_standards", seed.CsvFile))
	err = cmd.Run()
	if err != nil {
		slog.Error("error running exec: " + err.Error())
		return
	}

	os.Remove(seed.CsvFile)
	os.Remove(seed.SchemaFile)
	return
}
