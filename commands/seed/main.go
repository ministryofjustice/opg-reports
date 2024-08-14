package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

var flagset = flag.NewFlagSet("database seeder", flag.ExitOnError)

var (
	db     = argument.New(flagset, "db", "", "database file path")
	schema = argument.New(flagset, "schema", "", "schema file path")
	csv    = argument.New(flagset, "csv", "", "csv file pattern")
	table  = argument.New(flagset, "table", "github_standards", "table name")
	n      = argument.NewInt(flagset, "n", 1000, "number to generate")
)

type generatorF func(num int) []string

// map of funcs that generate fake data for seeding
var fakes map[string]generatorF = map[string]generatorF{
	"github_standards": func(num int) (lines []string) {
		lines = []string{}
		owner := fake.String(12)
		for x := 0; x < num; x++ {
			id := num + x
			g := ghs.Fake(&id, &owner)
			if x == 0 {
				lines = append(lines, g.CSVHead())
			}
			lines = append(lines, g.ToCSV())
		}
		return
	},
}

func Seed(dbV string, schemaV string, csvV string, tableV string, N int) (db *sql.DB, err error) {
	logger.LogSetup()
	slog.Info("starting to seed ",
		slog.String("db", dbV),
		slog.String("schema", schemaV),
		slog.String("csv", csvV),
		slog.String("table", tableV),
		slog.Int("n", N))
	if dbV == "" || schemaV == "" || csvV == "" || tableV == "" || N == 0 {
		err = fmt.Errorf("missing required arguments")
		return
	}
	// if the database exists, ignore
	dbExists := exists.FileOrDir(dbV)
	if dbExists {
		err = fmt.Errorf("database already exists: %s", dbV)
		return
	}
	// -- generate connection
	db, err = DB(dbV)
	if db != nil {
		defer db.Close()
	}

	if err != nil {
		return
	}

	// if the schema does not exist, error
	schemaExists := exists.FileOrDir(schemaV)
	if !schemaExists {
		err = fmt.Errorf("schema does not exist: %s", schemaV)
		return
	}
	// -- load schema into db
	err = SchemaLoad(db, schemaV)
	if err != nil {
		return
	}

	// -- now look for files matching the source pattern
	files := []string{}
	files, err = filepath.Glob(csvV)
	if err != nil {
		slog.Error("error with pattern matching", slog.String("csv", csvV), slog.String("err", err.Error()))
	}

	slog.Info("how many csv files found?", slog.Int("count", len(files)))

	// -- if we didnt find any files, then we try to make some fake ones
	if fk, ok := fakes[tableV]; ok && len(files) == 0 {
		csvDir := filepath.Dir(csvV)
		name := filepath.Join(csvDir, "generated_"+tableV+".csv")
		lines := fk(N)
		slog.Info("no files, but generator found so creating dummy data", slog.String("csv", name), slog.Int("line_count", len(lines)))
		f, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		defer f.Close()
		for _, line := range lines {
			f.WriteString(line)
		}
		files = append(files, name)
	}

	// -- now load files into the database
	for _, file := range files {
		slog.Info("data source", slog.String("csv", file))
		if exists.FileOrDir(file) {

			cmd := exec.Command("sqlite3", dbV, "--csv", fmt.Sprintf(".import --skip 1 %s %s", file, tableV))

			slog.Info("importing csv", slog.String("csv", file), slog.String("cmd", cmd.String()))
			err = cmd.Run()
			if err != nil {
				return
			}
			// remove the file
			// defer os.Remove(file)
		}
	}
	return
}

func SchemaLoad(db *sql.DB, schemaFile string) (err error) {
	schemaDir := filepath.Dir(schemaFile)
	os.MkdirAll(schemaDir, os.ModePerm)
	// -- read file
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return
	}
	return
}

// create a db
func DB(dbFile string) (db *sql.DB, err error) {
	dbDir := filepath.Dir(dbFile)
	os.MkdirAll(dbDir, os.ModePerm)
	// -- create a stub database file
	os.WriteFile(dbFile, []byte(""), os.ModePerm)

	// -- connect to the database and then import schema
	conn := consts.SQL_CONNECTION_PARAMS
	db, err = sql.Open("sqlite3", dbFile+conn)
	if err != nil && db != nil {
		db.Close()
	}
	return
}

func main() {
	flagset.Parse(os.Args[1:])

	// map args
	var (
		dbV     string = *db.Value
		schemaV string = *schema.Value
		csvV    string = *csv.Value
		tableV  string = *table.Value
		N       int    = *n.Value
	)

	if _, err := Seed(dbV, schemaV, csvV, tableV, N); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
