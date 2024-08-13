package seeder

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/exists"
)

func New(s *Seed) (db *sql.DB, err error) {

	// -- create directories
	schemaDir := filepath.Dir(s.Schema)
	dbDir := filepath.Dir(s.DB)
	os.MkdirAll(schemaDir, os.ModePerm)
	os.MkdirAll(dbDir, os.ModePerm)

	// -- create a stub database file
	os.WriteFile(s.DB, []byte(""), os.ModePerm)

	// -- connect to the database and then import schema
	conn := consts.SQL_CONNECTION_PARAMS
	db, err = sql.Open("sqlite3", s.DB+conn)
	if err != nil {
		db.Close()
		return
	}
	// -- read schema
	schema, err := os.ReadFile(s.Schema)
	if err != nil {
		db.Close()
		return
	}
	// -- load schema
	_, err = db.Exec(string(schema))
	if err != nil {
		db.Close()
		return
	}

	// -- if theres not set csv files, but we have lines, then make a new csv
	lines := s.Dummy
	if len(s.Source) == 0 && len(lines) > 0 {
		name := filepath.Join(dbDir, "generated.csv")
		slog.Info("generating csv from lines", slog.String("file", name))
		f, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		defer f.Close()
		for _, line := range lines {
			f.WriteString(line)
		}
		// add the created file into the stack to import
		s.Source = append(s.Source, name)
	}

	// -- load in the csv files
	for _, csv := range s.Source {
		slog.Info("data source", slog.String("csv", csv))
		if exists.FileOrDir(csv) {
			slog.Info("importing csv", slog.String("csv", csv))

			cmd := exec.Command("sqlite3", s.DB, "--csv", fmt.Sprintf(".import --skip 1 %s %s", csv, s.Table))
			err = cmd.Run()
			if err != nil {
				db.Close()
				return
			}
			// remove the file
			defer os.Remove(csv)
		}
	}

	return
}
