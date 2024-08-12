package main

import (
	"flag"
	"fmt"
	"log/slog"
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
	slog.Info("Seeding tables")
	// only generate data if it doesnt already exist
	if what == "github_standards" || what == "all" || what == "" {
		d := *dir.Value
		dbDir := fmt.Sprintf("%s/dbs", d)

		if !exists.FileOrDir(dbDir + "/github_standards.db") {
			slog.Info("Seeding github_standards")
			github_standards_seed.NewSeed(dbDir, 1000)
			slog.Info("Seeded github_standards")
		} else {
			slog.Info("github_standards exists")
		}
	}

}
