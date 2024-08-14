package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
)

var flagset = flag.NewFlagSet("database seeder", flag.ExitOnError)

var (
	db     = argument.New(flagset, "db", "", "database file path")
	schema = argument.New(flagset, "schema", "", "schema file path")
	data   = argument.New(flagset, "data", "", "data file pattern")
	table  = argument.New(flagset, "table", "github_standards", "table name")
	n      = argument.NewInt(flagset, "n", 1000, "number to generate")
)

func main() {
	ctx := context.Background()
	flagset.Parse(os.Args[1:])
	// map args
	var (
		dbV     string = *db.Value
		schemaV string = *schema.Value
		dataV   string = *data.Value
		tableV  string = *table.Value
		N       int    = *n.Value
	)

	if _, err := seeder.Seed(ctx, dbV, schemaV, dataV, tableV, N); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
