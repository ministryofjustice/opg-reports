package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/seeder/github_standards_seed"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func main() {
	logger.LogSetup()
	group := flag.NewFlagSet("demo", flag.ExitOnError)
	which := argument.New(group, "which", "all", "")
	dir := argument.New(group, "dir", ".", "")

	group.Parse(os.Args[1:])

	what := *which.Value
	ctx := context.Background()
	// only generate data if it doesnt already exist
	if what == "github_standards" || what == "all" {
		d := *dir.Value
		dbDir := fmt.Sprintf("%s/dbs", d)
		schema := d + "/github_standards.sql"
		dbPath := dbDir + "/github_standards.db"
		if !exists.FileOrDir(dbPath) {
			os.MkdirAll(d, os.ModePerm)
			db := github_standards_seed.NewDb(ctx, dbPath, schema)
			github_standards_seed.Seed(ctx, db, 100)
		}
	}

}
