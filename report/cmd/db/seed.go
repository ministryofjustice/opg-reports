package main

import (
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	seedCmdName   string = "seed" // root command name
	seedShortDesc string = `seed inserts test data into the configured database`
	seedLongDesc  string = `
seed inserts test data into the configured database; generally intended for development use only. This will also run database migrataions before inserting seed data.
`
)

var seedCmd = &cobra.Command{
	Use:   seedCmdName,
	Short: seedShortDesc,
	Long:  seedLongDesc,
	RunE:  seedRunE,
}

var (
	seedDBPath   string = "database/api.db" // represents --db
	seedDBDriver string = "sqlite3"         // represents --driver
)

func seedRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	// db connection
	db, err = dbconnection.Connection(ctx, log, seedDBDriver, seedDBPath)
	if err != nil {
		return
	}
	// close the db
	defer db.Close()

	err = dbsetup.SeedAll(ctx, log, db)
	if err != nil {
		return
	}

	return
}

// setup default values for config and logging
func init() {
	seedCmd.Flags().StringVar(&seedDBPath, "db", seedDBPath, "Path to database")
	seedCmd.Flags().StringVar(&seedDBDriver, "driver", seedDBDriver, "Database driver")
}
